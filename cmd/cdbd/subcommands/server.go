package subcommands

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/config-source/cdb/internal/configvalues"
	"github.com/config-source/cdb/internal/repository/postgres"
	"github.com/config-source/cdb/internal/server"
	"github.com/config-source/cdb/internal/server/middleware"
	"github.com/config-source/cdb/internal/settings"
	"github.com/pseidemann/finish"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func setupLogger() zerolog.Logger {
	logger := zerolog.New(os.Stdout).
		Level(settings.LogLevel()).
		With().
		Timestamp().
		Logger()
	if settings.HumanLogs() {
		logger = logger.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}
	// Set durations to render as Milliseconds
	zerolog.DurationFieldUnit = time.Millisecond
	return logger
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := setupLogger()

		repo, err := postgres.NewRepository(
			context.Background(),
			logger,
			settings.DBUrl(),
		)
		if err != nil {
			return err
		}

		var server http.Handler = server.New(
			repo,
			configvalues.NewService(repo, settings.DynamicConfigKeys()),
			logger,
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
