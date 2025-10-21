package mylogger

import (
	"io"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)



type LogConfig struct {
	ConsoleLoggingEnabled bool
	FileLoggingEnabled    bool
	Directory             string
	Filename              string
	MaxSizeMB             int
	MaxBackups            int
	MaxAgeDays            int
}

func ConfigureLogger(config LogConfig) zerolog.Logger {
	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
		writers = append(writers, consoleWriter)
	}

	if config.FileLoggingEnabled {
		if err := os.MkdirAll(config.Directory, 0744); err != nil {
			log.Error().Err(err).Str("path", config.Directory).Msg("Failed to create log directory")
		} else {
			fileLogger := &lumberjack.Logger{
				Filename:   path.Join(config.Directory, config.Filename),
				MaxBackups: config.MaxBackups,
				MaxSize:    config.MaxSizeMB,
				MaxAge:     config.MaxAgeDays,
				Compress:   true,
			}
			writers = append(writers, fileLogger)
		}
	}

	mw := io.MultiWriter(writers...)

	newLogger := zerolog.New(mw).With().
		Timestamp().
		Logger()

	return newLogger.Level(zerolog.InfoLevel)
}


func HandleLogging(logger zerolog.Logger,err error,message string){
	serviceLogger := logger .With().Caller().Logger()
		serviceLogger.Error().Err(err).Str("userID", "userID").Msg(message)
}