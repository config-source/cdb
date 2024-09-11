package commands

import (
	"context"
	"fmt"
	"syscall"

	"github.com/config-source/cdb/cmd/cdb/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func login(ctx context.Context) (string, error) {
	email, err := input("Email: ")
	if err != nil {
		return "", err
	}

	fmt.Print("Password: ")
	bytepw, err := term.ReadPassword(syscall.Stdin)
	fmt.Print("\n")
	if err != nil {
		return "", err
	}

	_, err = config.Client.Login(ctx, email, string(bytepw))
	if err != nil {
		return "", err
	}

	token, err := config.Client.IssueAPIToken(ctx, email, string(bytepw))
	return token.Token, err
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Generate an API token for use with this CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := login(cmd.Context())
		if err != nil {
			return err
		}

		fmt.Println(token)
		return nil
	},
}
