package main

import (
	"errors"
	"net/http"

	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/auth/postgres"
	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/configvalues"
	"github.com/config-source/cdb/environments"
	"github.com/config-source/cdb/server"
	"github.com/config-source/cdb/server/middleware"
	"github.com/config-source/cdb/services"
	"github.com/config-source/cdb/settings"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pseidemann/finish"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func getAuthenticationGateway(log zerolog.Logger, pool *pgxpool.Pool) auth.AuthenticationGateway {
	gatewayName := settings.AuthenticationGateway()
	switch gatewayName {
	default:
		if gatewayName == "" {
			log.Warn().Msg("no AUTHENTICATION_GATEWAY configured, using postgres as default")
		} else if gatewayName != "postgres" {
			log.Error().Str("gatewayName", gatewayName).Msg("is not a valid gateway")
		}

		return postgres.NewGateway(log, pool)
	}
}

func getAuthorizationGateway(log zerolog.Logger, pool *pgxpool.Pool) auth.AuthorizationGateway {
	gatewayName := settings.AuthorizationGateway()
	switch gatewayName {
	default:
		if gatewayName == "" {
			log.Warn().Msg("no AUTHORIZATION_GATEWAY configured, using postgres as default")
		} else if gatewayName != "postgres" {
			log.Error().Str("gatewayName", gatewayName).Msg("is not a valid gateway")
		}

		return postgres.NewGateway(log, pool)
	}
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := settings.GetLogger()

		pool, err := pgxpool.New(
			cmd.Context(),
			settings.DBUrl(),
		)
		if err != nil {
			return err
		}

		authenticationGateway := getAuthenticationGateway(logger, pool)
		authorizationGateway := getAuthorizationGateway(logger, pool)

		envsRepo := environments.NewRepository(logger, pool)
		keysRepo := configkeys.NewRepository(logger, pool)
		valuesRepo := configvalues.NewRepository(logger, pool)
		svcRepo := services.NewRepository(logger, pool)
		tokenRegistry := auth.NewTokenRegistry(logger, pool)

		envsService := environments.NewService(envsRepo, authorizationGateway)
		keysService := configkeys.NewService(keysRepo, authorizationGateway)
		valuesService := configvalues.NewService(
			valuesRepo,
			envsRepo,
			keysRepo,
			authorizationGateway,
			settings.DynamicConfigKeys(),
		)
		userService := auth.NewUserService(
			authenticationGateway,
			authorizationGateway,
			tokenRegistry,
			settings.AllowPublicRegistration(),
			settings.DefaultRegisterRole(),
		)
		svcService := services.NewServiceService(svcRepo, authorizationGateway)

		var server http.Handler = server.New(
			logger,
			settings.JWTSigningKey(),
			pool,
			userService,
			valuesService,
			envsService,
			keysService,
			svcService,
			settings.FrontendLocation(),
		)
		server = middleware.AccessLog(logger, server)

		httpServer := &http.Server{Addr: settings.ListenAddr(), Handler: server}

		fin := finish.New()
		fin.Add(httpServer)

		go func() {
			if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
				logger.Err(err).Msg("error closing down http server")
			}
		}()

		logger.Info().Str("address", settings.ListenAddr()).Msg("listening for connections")
		fin.Wait()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
