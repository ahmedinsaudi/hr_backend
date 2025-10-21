package myfiber

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/jmoiron/sqlx"
	"githup.ahmedramadan.4cashier/internal/bootstrap"
	"githup.ahmedramadan.4cashier/internal/handler"
)

// setup fiber
func SetupFiber(handlers Handlers, db *sqlx.DB) {
	app := fiber.New(fiber.Config{
		AppName:       "HADEF Fiber App",
		ReadTimeout:   10 * time.Second,
		WriteTimeout:  10 * time.Second,
		IdleTimeout:   60 * time.Second,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
	})

	app.Use(cors.New(cors.Config{}))
	app.Use(logger.New())

	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
	}))

	app.Use("/", static.New("./public"))

	//handler.JWTAuthMiddleware(), handler.HasRolesMiddleware("employee", "admin")

	hrGroup := app.Group("/hr")

	hrGroup.Get("/hr-profiles", handlers.HRHandler.GetHRProfiles)     // Get HR Profiles
	hrGroup.Get("/rates", handlers.HRHandler.GetRates)                // Get HR rates
	hrGroup.Post("/:id/experience", handlers.HRHandler.AddExperience) // Add experience to HR
	hrGroup.Post("/:id/job-roles", handlers.HRHandler.AddJobRoles)    // Add job roles to HR
	hrGroup.Post("/rate", handlers.HRHandler.RateHR)                  // Rate an HR profile
	hrGroup.Post("/rate/like", handlers.HRHandler.LikeRate)           // Like a HR rate
	hrGroup.Post("/badge", handlers.HRHandler.AwardBadge)             // Award a badge to HR
	hrGroup.Post("/badge/like", handlers.HRHandler.LikeBadge)         // Like a badge for HR
	hrGroup.Get("/:employee_id/stats", handlers.HRHandler.GetEmployeeStats)

	app.Post("/signin", handler.SignIn(db))
	app.Post("/signup", handler.SignUpUser(db))

	admin := app.Group("/api/admin", handler.JWTAuthMiddleware(), handler.HasRolesMiddleware("admin"))
	admin.Get("/dashboard", bootstrap.AdminEndpoint)

	app.Get("/zat", func(c fiber.Ctx) error {

		return c.JSON(fiber.Map{"message": "Welcome, 55 Editor! Here is your content."})
	})

	app.Get("/health", func(c fiber.Ctx) error {

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	editor := app.Group("/api/editor", handler.JWTAuthMiddleware(), handler.HasRolesMiddleware("editor"))

	editor.Get("/content", func(c fiber.Ctx) error {

		userClaims, ok := c.Locals("user").(*bootstrap.UserClaims)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user claims from context"})
		}
		return c.JSON(fiber.Map{"message": "Welcome, Editor! Here is your content.", "user_id": userClaims.ID})
	})

	log.Fatal(app.Listen(":8080"))
}
