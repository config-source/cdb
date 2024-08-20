package subcommands

import (
	"context"
	"fmt"
	"syscall"

	"github.com/config-source/cdb/internal/auth/postgres"
	"github.com/config-source/cdb/internal/settings"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var createAdminUserCmd = &cobra.Command{
	Use:   "create-admin-user",
	Short: "Create an Administrator user",
	Long: `Create a CDB Administrator user.

The created user will have FULL access to CDB including managing users and
roles. This command is useful when you've locked yourself out or are setting
up the initial management user.

Defaults to assigning the Administrator role that comes with CDB. If you've
changed that role or want to use a different role then specify the --role
flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := settings.GetLogger()

		gateway, err := postgres.NewGateway(
			context.Background(),
			logger,
			settings.DBUrl(),
		)
		if err != nil {
			return err
		}

		fmt.Print("Password: ")
		bytepw, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			return err
		}

		user, err := gateway.Register(context.Background(), email, string(bytepw))
		if err != nil {
			return err
		}

		return gateway.AssignRoleToUserNoAuth(context.Background(), user, role)
	},
}

var email string
var role string

func init() {
	createAdminUserCmd.Flags().StringVarP(&email, "email", "e", "", "")
	createAdminUserCmd.MarkFlagRequired("email")

	createAdminUserCmd.Flags().StringVarP(&role, "role", "r", "Administrator", "")

	rootCmd.AddCommand(createAdminUserCmd)
}
