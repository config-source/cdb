package settings

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/auth/postgres"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func init() {
	logger = zerolog.New(os.Stdout).
		Level(LogLevel()).
		With().
		Timestamp().
		Logger()
	if HumanLogs() {
		logger = logger.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}
	// Set durations to render as Milliseconds
	zerolog.DurationFieldUnit = time.Millisecond

	auth.TokenIssuer = TokenIssuer()
}

// TokenIssuer returns the JWT token issuer
func TokenIssuer() string {
	issuer := os.Getenv("JWT_TOKEN_ISSUER")
	if issuer == "" {
		return "cdb"
	}

	return issuer
}

// GetLogger returns the preconfigured global logger
func GetLogger() zerolog.Logger {
	return logger
}

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
	val := os.Getenv("DYNAMIC_CONFIG_KEYS")
	if val == "" {
		val = "true"
	}

	return strings.ToLower(val) == "true"
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
		return "0.0.0.0:8080"
	}

	return addr
}

// FrontendLocation returns the location of the statically built frontend files
// or an http upstream to proxy to.
func FrontendLocation() string {
	location := os.Getenv("FRONTEND_LOCATION")
	if location == "" {
		return "./frontend/build"
	}

	return location
}

var keyCache []byte

// JWTSigningKey returns the configured symmetric key for signing JWTs issued by
// CDB's auth system.
func JWTSigningKey() []byte {
	if keyCache == nil {
		key := os.Getenv("JWT_SIGNING_KEY")
		if key == "" {
			logger := GetLogger()
			logger.
				Error().
				Msg(
					"JWT_SIGNING_KEY must be configured",
				)
			os.Exit(1)
		}

		keyCache = []byte(key)
	}

	return keyCache
}

func GetAuthenticationGateway(ctx context.Context, log zerolog.Logger) auth.AuthenticationGateway {
	gatewayName := os.Getenv("AUTHENTICATION_GATEWAY")
	switch gatewayName {
	default:
		if gatewayName == "" {
			log.Warn().Msg("no AUTHENTICATION_GATEWAY configured, using postgres as default")
		}

		gw, err := postgres.NewGateway(ctx, log, DBUrl())
		if err != nil {
			log.Panic().Err(err).Msg("Failed to load gateway")
		}

		return gw
	}
}

func GetAuthorizationGateway(ctx context.Context, log zerolog.Logger) auth.AuthorizationGateway {
	gatewayName := os.Getenv("AUTHORIZATION_GATEWAY")
	switch gatewayName {
	default:
		if gatewayName == "" {
			log.Warn().Msg("no AUTHORIZATION_GATEWAY configured, using postgres as default")
		}

		gw, err := postgres.NewGateway(ctx, log, DBUrl())
		if err != nil {
			log.Panic().Err(err).Msg("Failed to load gateway")
		}

		return gw
	}
}

func AllowPublicRegistration() bool {
	val := os.Getenv("ALLOW_PUBLIC_REGISTRATION")
	return strings.ToLower(val) == "true"
}

func DefaultRegisterRole() string {
	val := os.Getenv("DEFAULT_REGISTER_ROLE")
	if !AllowPublicRegistration() {
		return ""
	} else if val != "" {
		logger.Warn().
			Msg("ALLOW_PUBLIC_REGISTRATION is on but DEFAULT_REGISTER_ROLE is empty so it will not work")
	}

	return val
}
