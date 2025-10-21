package bootstrap

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

const jwtSecret = "your-super-secret-jwt-key"


type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Role    string `json:"role"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type UserClaims struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}



func AdminEndpoint(c fiber.Ctx) error {
	userClaims, ok := c.Locals("user").(*UserClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user claims from context"})
	}

	return c.JSON(fiber.Map{
		"message":    "Welcome, Admin!",
		"user_id":    userClaims.ID,
		"user_email": userClaims.Email,
	})
}



func SeedAdminDirect(tx *sqlx.Tx, branchID int, companyID int) error {
    // Generate password hash
    hash, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    // Current timestamp
    now := time.Now()

    // Admin email based on branch
    adminEmail := fmt.Sprintf("d%d@d", branchID)

    var userID int
    err = tx.QueryRowx(`
        INSERT INTO users 
            (name, email, password_hash, user_type, is_active, created_at, updated_at)
        VALUES 
            ('admin', $1, $2, 'employee', true, $3, $3)
        RETURNING id
    `, adminEmail, string(hash), now).Scan(&userID)

    if err != nil {
        return fmt.Errorf("failed to insert admin user: %w", err)
    }

    // Insert employee record
    _, err = tx.Exec(`
        INSERT INTO employees (user_id, branch_id, company_id, role)
        VALUES ($1, $2, $3, 'admin')
    `, userID, branchID, companyID)
    if err != nil {
        return fmt.Errorf("failed to insert admin employee record: %w", err)
    }

    return nil
}