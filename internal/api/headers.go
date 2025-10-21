package myfiber

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"githup.ahmedramadan.4cashier/internal/bootstrap"
	"githup.ahmedramadan.4cashier/internal/handler"
	"githup.ahmedramadan.4cashier/internal/repos"
	"githup.ahmedramadan.4cashier/internal/service"
)

type Handlers struct {
	
	HRHandler             handler.HRHandler
	EmployeeHandler             handler.EmployeeHandler
}

type App struct {
	DB       *sqlx.DB
	Handlers Handlers
}

func SetupHandlers(logger zerolog.Logger) *App {
	db := bootstrap.InitDB()

	// accountRepo := repos.NewPosAccountRepository(db)
	// accountService := service.NewAccountService(logger, accountRepo)
	// accountHandler := handler.NewAccountHandler(logger, accountService)



	employeeRepo := repos.NewPosEmployeeRepository(db)
	employeeService := service.NewEmployeeService(logger, employeeRepo)
	employeeHandler := handler.NewEmployeeHandler(logger, employeeService)

	hrRepo := repos.NewPosHRRepository(db)
	hrService := service.NewHRService(logger, hrRepo)
	hrHandler := handler.NewHRHandler(logger, hrService)

	return &App{
		DB: db,
		Handlers: Handlers{
			HRHandler:              *hrHandler,
			EmployeeHandler:            *employeeHandler,
		},
	}
}
