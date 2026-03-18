package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Log zerolog.Logger

func Init() {
	env := os.Getenv("ENV")
	
	if env == "production" {
		// Production: JSON format para Cloud Logging
		zerolog.TimeFieldFormat = time.RFC3339
		zerolog.LevelFieldName = "severity"
		zerolog.TimestampFieldName = "timestamp"
		
		Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		// Development: Console output colorido e legível
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		Log = zerolog.New(output).With().Timestamp().Logger()
	}
	
	// Set global logger
	log.Logger = Log
	
	Log.Info().Str("env", env).Msg("Logger initialized")
}

func GetLogger() zerolog.Logger {
	return Log
}
