package settings

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// DBUrl returns the connection string as set by $DB_URL
//
// Most often this is left empty and instead the PG* variables are used for
// connecting to the database. This is the recommended approach.
func DBUrl() string {
	return os.Getenv("DB_URL")
}

// DynamicConfigKeys indicates that Config Keys should be created if they aren't
// found when a new Config Value is created.
func DynamicConfigKeys() bool {
	return strings.ToLower(os.Getenv("DYNAMIC_CONFIG_KEYS")) == "true"
}

// HumanLogs indicates that structured logging should not be used and instead
// logs should be human-friendly.
func HumanLogs() bool {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	return env == "local"
}

// LogLevel returns the configured log level as defined by $LOG_LEVEL
//
// Defaults to InfoLevel
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

// LogLevel returns the configured listen address as defined by $LISTEN_ADDRESS
//
// Defaults to all interfaces listening on port 8080.
func ListenAddr() string {
	addr := os.Getenv("LISTEN_ADDRESS")
	if addr == "" {
		return ":8080"
	}

	return addr
}
