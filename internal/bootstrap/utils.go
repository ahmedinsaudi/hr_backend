package bootstrap

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error in load .env :%v", err)
	}
}

func InitDB() *sqlx.DB {
	dns := os.Getenv("DATABASE_URL")
	if dns == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sqlx.Connect("postgres", dns)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	const maxOpenConns = 4
	db.SetMaxOpenConns(maxOpenConns)

	const maxIdleConns = 2
	db.SetMaxIdleConns(maxIdleConns)

	db.SetConnMaxIdleTime(2 * time.Second)

	db.SetConnMaxLifetime(1 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("error pinging database: %v", err)
	}

	log.Printf("Successfully connected to database with pool: MaxOpen=%d, MaxIdle=%d", maxOpenConns, maxIdleConns)
	return db
}

type Pagination struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func GetPagination(ctx fiber.Ctx) Pagination {
	pagination := Pagination{
		Page:  1,  // Default page number
		Limit: 10, // Default items per page
	}

	// Try to parse 'page' parameter
	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			pagination.Page = page
		}
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			pagination.Limit = limit
		}
	}

	return pagination
}

func BindGenericRequestBody[T any](targetType T) fiber.Handler {
	return func(c fiber.Ctx) error {
		val := reflect.New(reflect.TypeOf(targetType)).Interface().(*T)

		if err := c.Bind().Body(val); err != nil {
			log.Printf("error is=%v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true})
		}

		c.Locals("bodyParse", val)

		return c.Next()
	}
}
