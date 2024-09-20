package env

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/config-source/cdb/cmd/cdb/config"
	"github.com/spf13/cobra"
)

var envGetCmd = &cobra.Command{
	Use:   "get <service-name> <environment-name>",
	Short: "Get environment information by name",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := config.Client.GetEnvironmentByName(context.Background(), args[0], args[1])
		if err != nil {
			return err
		}

		output, err := json.MarshalIndent(env, "", "    ")
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	Command.AddCommand(envGetCmd)
}
