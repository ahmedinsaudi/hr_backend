FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o /hr_rate .

FROM alpine:latest
RUN apk add --no-cache tzdata ca-certificates bash wget
RUN addgroup -S -g 1000 app && adduser -S -u 1000 -G app app
RUN mkdir -p /app/application_logs && chown -R app:app /app/application_logs

COPY --from=builder /hr_rate /usr/local/bin/
USER app
WORKDIR /app
EXPOSE 7070
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s CMD wget -qO- http://127.0.0.1:7070/health || exit 1
ENTRYPOINT ["/usr/local/bin/hr_rate"]
