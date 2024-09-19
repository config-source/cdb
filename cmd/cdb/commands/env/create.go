package env

import (
	"context"
	"errors"
	"fmt"

	"github.com/config-source/cdb/cmd/cdb/config"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/spf13/cobra"
)

var promotesTo string

func getPromotesToID() int {
	env, err := config.Client.GetEnvironmentByNameOrID(promotesTo)
	if err == nil {
		return env.ID
	}

	return 0
}

var envCreateCmd = &cobra.Command{
	Use:   "create <environment-name>",
	Short: "Create a new environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		env := environments.Environment{
			Name: args[0],
		}

		if promotesTo != "" {
			promotesToID := getPromotesToID()
			if promotesToID == 0 {
				return errors.New("unable to identify the environment you want to promote to, is the name correct?")
			}

			env.PromotesToID = &promotesToID
		}

		_, err := config.Client.CreateEnvironment(context.Background(), env)
		if err != nil {
			return err
		}

		fmt.Println("Successfully created environment:", env.Name)
		return nil
	},
}

func init() {
	envCreateCmd.Flags().StringVarP(&promotesTo, "promotes-to", "p", "", "What environment this promotes to, accepts an environment name or ID.")
	Command.AddCommand(envCreateCmd)
}
