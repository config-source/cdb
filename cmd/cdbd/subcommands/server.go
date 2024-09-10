package subcommands

import (
	"errors"
	"net/http"

	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/repository/postgres"
	"github.com/config-source/cdb/server"
	"github.com/config-source/cdb/server/middleware"
	"github.com/config-source/cdb/services"
	"github.com/config-source/cdb/settings"
	"github.com/pseidemann/finish"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := settings.GetLogger()

		repo, err := postgres.NewRepository(
			cmd.Context(),
			logger,
			settings.DBUrl(),
		)
		if err != nil {
			return err
		}

		authenticationGateway := settings.GetAuthenticationGateway(cmd.Context(), logger)
		authorizationGateway := settings.GetAuthorizationGateway(cmd.Context(), logger)

		var server http.Handler = server.New(
			repo,
			logger,
			settings.JWTSigningKey(),
			auth.NewUserService(
				authenticationGateway,
				authorizationGateway,
				settings.AllowPublicRegistration(),
				settings.DefaultRegisterRole(),
			),
			services.NewConfigValuesService(
				repo,
				authorizationGateway,
				settings.DynamicConfigKeys(),
			),
			services.NewEnvironmentsService(repo, authorizationGateway),
			services.NewConfigKeysService(repo, authorizationGateway),
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
