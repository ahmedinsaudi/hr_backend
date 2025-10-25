package main

import (
	myfiber "githup.ahmedramadan.4cashier/internal/api"
	// "githup.ahmedramadan.4cashier/internal/bootstrap"
	mylogger "githup.ahmedramadan.4cashier/internal/mylogger"
)

func main() {
	// bootstrap.LoadEnv()

	logConfig := mylogger.LogConfig{
		ConsoleLoggingEnabled: true,               // Enable logging to console
		FileLoggingEnabled:    true,               // Enable logging to file
		Directory:             "application_logs", // Directory where log files will be stored
		Filename:              "app.log",          // Base name for the log file
		MaxSizeMB:             10,                 // Max size of the log file before rotation (10 MB)
		MaxBackups:            5,                  // Keep up to 5 old log files
		MaxAgeDays:            7,                  // Delete log files older than 7 days
	}

	appLogger := mylogger.ConfigureLogger(logConfig)
	appLogger.Info().Msg("Application starting...")

	app := myfiber.SetupHandlers(appLogger)

	myfiber.SetupFiber(app.Handlers, app.DB)
}

// bootstrap

// Config
