package handler

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"githup.ahmedramadan.4cashier/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var jwtSecret = []byte("supersecretkey") // change in production

type UserClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	JobField     string `json:"job_field" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	UserType string `json:"userType" validate:"required,oneof=hr employee"`
}

func SignIn(db *sqlx.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req LoginRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		var (
			role   string
			userID int
			email  string
		)

		var emp models.Employee
		err := db.Get(&emp, "SELECT * FROM employees WHERE email = $1", req.Email)
		log.Printf("paswprd %v",emp)
		if err == nil && emp.PasswordHash != "" {

			if err := bcrypt.CompareHashAndPassword([]byte(emp.PasswordHash), []byte(req.Password)); err != nil {
				log.Print(err)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
			}
			userID = emp.ID
			email = emp.Email
			role = "employee"
		} else {
			var hr models.HRProfile
			err := db.Get(&hr, "SELECT * FROM hr_profiles WHERE email = $1", req.Email)
			if err != nil || hr.PasswordHash == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
			}
			if err := bcrypt.CompareHashAndPassword([]byte(hr.PasswordHash), []byte(req.Password)); err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
			}
			userID = hr.ID
			email = *hr.Email
			role = "hr"
		}

		claims := UserClaims{
			UserID: userID,
			Email:  email,
			Role:   role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create token"})
		}

		return c.JSON(fiber.Map{
			"token": tokenString,
			"user": fiber.Map{
				"id":    userID,
				"email": email,
				"role":  role,
			},
		})
	}
}

func JWTAuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid Authorization header"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &UserClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// ✅ store claims for later use
		c.Locals("user", claims)
		return c.Next()
	}
}

func HasRolesMiddleware(requiredRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userClaims, ok := c.Locals("user").(*UserClaims)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "User claims not found"})
		}

		for _, requiredRole := range requiredRoles {
			if userClaims.Role == requiredRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}
}

func SignUpUser(db *sqlx.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req CreateUserRequest

		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		if err := models.Validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
		}

		now := time.Now()
		tx, err := db.Beginx()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
		}
		defer func() {
			if p := recover(); p != nil {
				tx.Rollback()
				panic(p)
			}
		}()

		switch strings.ToLower(req.UserType) {
		case "employee":

			employee := models.Employee{
				Name:         req.Name,
				Email:        req.Email,
				PasswordHash: string(hash),
				JobField: req.JobField,
				IsVerified:   false,
				CreatedAt:    now,
				UpdatedAt:    now,
			}

			query := `
				INSERT INTO employees (name, email, password_hash,job_field, is_verified, created_at, updated_at)
				VALUES (:name, :email, :password_hash,:job_field, :is_verified, :created_at, :updated_at)
				RETURNING id
			`
			rows, err := tx.NamedQuery(query, &employee)
			log.Print(err)
			if err != nil {
				
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert employee: " + err.Error()})
			}
			if rows.Next() {
				if err := rows.Scan(&employee.ID); err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get employee ID"})
				}
			}

			if err := tx.Commit(); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
			}
			
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Employee account created", "id": employee.ID})

		case "hr":
			hr := models.HRProfile{
				Name:         &req.Name,
				Email:        &req.Email,
				PasswordHash: string(hash),
				JobPosition: &req.JobField,
				Verified:     false,
				CreatedAt:    now,
				UpdatedAt:    now,
			}

			query := `
				INSERT INTO hr_profiles (name, email, password_hash, verified_profile, created_at, updated_at)
				VALUES (:name, :email, :password_hash, :verified_profile, :created_at, :updated_at)
				RETURNING id
			`
			rows, err := tx.NamedQuery(query, &hr)
			if err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert HR profile: " + err.Error()})
			}
			if rows.Next() {
				if err := rows.Scan(&hr.ID); err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get HR ID"})
				}
			}

			if err := tx.Commit(); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
			}

			return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "HR account created", "id": hr.ID})

		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user type"})
		}
	}
}
