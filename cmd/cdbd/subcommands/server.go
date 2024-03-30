package subcommands

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/config-source/cdb/internal/postgres"
	"github.com/config-source/cdb/internal/server/api"
	"github.com/config-source/cdb/internal/settings"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		repo, err := postgres.NewRepository(
			context.Background(),
			logger,
			settings.DBUrl(),
		)
		if err != nil {
			return err
		}

		server := api.New(repo, logger)
		return http.ListenAndServe(":8080", server)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}