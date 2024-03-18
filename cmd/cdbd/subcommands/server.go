package subcommands

import (
	"context"
	"net/http"

	"github.com/config-source/cdb/internal/postgres"
	"github.com/config-source/cdb/internal/server/api"
	"github.com/config-source/cdb/internal/settings"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := postgres.NewRepository(context.Background(), settings.DBUrl())
		if err != nil {
			return err
		}

		server := api.New(repo)
		return http.ListenAndServe(":8080", server)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
