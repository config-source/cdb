package main

import (
	"net/http"

	"github.com/config-source/cdb/internal/settings"
	"github.com/config-source/cdb/pkg/configvalues"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/webhooks"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

var testWebhookCmd = &cobra.Command{
	Use:   "test-webhook",
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

		repo := webhooks.NewRepository(logger, pool)
		envsRepo := environments.NewRepository(logger, pool)
		valuesRepo := configvalues.NewRepository(logger, pool)

		env, err := envsRepo.GetEnvironmentByName(cmd.Context(), "dev")
		if err != nil {
			return err
		}

		whDefs, err := repo.GetWebhookDefinitionsForEnvironment(cmd.Context(), env.ID)
		if err != nil {
			return err
		}

		config, err := valuesRepo.GetConfiguration(cmd.Context(), env.ID)
		if err != nil {
			return err
		}

		for _, wh := range whDefs {
			err = webhooks.RunWebhook(cmd.Context(), logger, http.DefaultClient, wh, env, config)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testWebhookCmd)
}
