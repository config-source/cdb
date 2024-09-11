package commands

import (
	"os"

	"github.com/config-source/cdb/cmd/cdb/commands/configuration"
	"github.com/config-source/cdb/cmd/cdb/commands/env"
	"github.com/config-source/cdb/cmd/cdb/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "cdb",
	Short:        "Command line interface for your configuration database.",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.LoadConfig()
	},
}

func init() {
	rootCmd.AddCommand(configuration.Command)
	rootCmd.AddCommand(env.Command)
	rootCmd.AddCommand(setupCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
