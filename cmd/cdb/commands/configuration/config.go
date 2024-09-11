package configuration

import "github.com/spf13/cobra"

var Command = &cobra.Command{
	Use: "configuration <subcommand>",
	Aliases: []string{
		"c",
		"cfg",
		"config",
	},
}

func init() {
	Command.AddCommand(getConfigCmd)
}
