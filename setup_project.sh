#!/bin/bash
set -e

if [ -z "$1" ]; then
echo "❌ استخدم: ./setup_project.sh "
exit 1
fi

PROJECT_NAME=$1
PROJECT_PORT=${2:-8080}
PROJECT_PATH="/srv/projects/${PROJECT_NAME}"
DB_NAME="${PROJECT_NAME}_db"

echo "🚀 إنشاء مشروع جديد: ${PROJECT_NAME} على بورت ${PROJECT_PORT}"
echo "📁 المجلد: ${PROJECT_PATH}"

# إنشاء مجلد المشروع

mkdir -p "${PROJECT_PATH}/migrations"
mkdir -p "${PROJECT_PATH}/application_logs"
mkdir -p "${PROJECT_PATH}/public/images"

# توليد ملف البيئة

cat > "${PROJECT_PATH}/.env" <<EOL
POSTGRES_USER=postgres
POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-changeme}
POSTGRES_DB=${DB_NAME}
DATABASE_URL=postgres://postgres:${POSTGRES_PASSWORD}@pgbouncer:5432/${DB_NAME}?sslmode=disable
JWT_SECRET=$(openssl rand -hex 32)
APP_ENV=production
EOL

echo "🧱 إنشاء قاعدة بيانات ${DB_NAME}..."
docker exec infra-db-1 psql -U postgres -c "CREATE DATABASE ${DB_NAME};" || true

# توليد ملف docker-compose للمشروع

cat > "${PROJECT_PATH}/docker-compose.yml" <<EOL
services:
${PROJECT_NAME}:
image: ghcr.io/${GHCR_USER}/${PROJECT_NAME}:latest
restart: unless-stopped
env_file: .env
expose:
- "${PROJECT_PORT}"
environment:
- VIRTUAL_HOST=${PROJECT_NAME}.doneally.com
- VIRTUAL_PORT=${PROJECT_PORT}
- VIRTUAL_PROTO=http
- LETSENCRYPT_HOST=${PROJECT_NAME}.doneally.com
- LETSENCRYPT_EMAIL=[admin@doneally.com](mailto:admin@doneally.com)
depends_on:
- pgbouncer
networks:
- backend
- frontend
volumes:
- ./migrations:/migrations
- ./application_logs:/app/application_logs
- ./public/images:/app/public/images

networks:
frontend:
external: true
backend:
external: true
EOL

echo "📦 تسجيل الدخول في GHCR..."
echo "${GHCR_LOGIN}" | docker login ghcr.io -u ${GHCR_USER} --password-stdin

echo "⬇️ سحب وتشغيل المشروع ${PROJECT_NAME}..."
cd "${PROJECT_PATH}"
docker compose pull
docker compose up -d

echo "🧠 تشغيل المايجريشن..."
docker compose run --rm ${PROJECT_NAME} 
sh -c 'export PGPASSWORD=$POSTGRES_PASSWORD; for f in /migrations/*.sql; do echo "Running $f..."; psql -h pgbouncer -U postgres -d $POSTGRES_DB -f "$f"; done'

echo "✅ تم تشغيل ${PROJECT_NAME} بنجاح على https://${PROJECT_NAME}.doneally.com"
