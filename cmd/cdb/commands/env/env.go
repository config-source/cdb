package env

import "github.com/spf13/cobra"

var Command = &cobra.Command{
	Use: "environment <subcommand>",
	Aliases: []string{
		"e",
		"env",
	},
}

func init() {
	Command.AddCommand(envGetCmd)
	Command.AddCommand(envTreeCmd)
}
