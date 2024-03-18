package settings

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func DBUrl() string {
	// Empty string will use the PG* variables
	return os.Getenv("DB_URL")
}

func HumanLogs() bool {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	return env == "local"
}

func LogLevel() zerolog.Level {
	logLevel := os.Getenv("LOG_LEVEL")
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		return zerolog.DebugLevel
	case "ERROR":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
