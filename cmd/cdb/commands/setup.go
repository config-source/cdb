package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/config-source/cdb/cmd/cdb/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func input(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	ans, err := reader.ReadString('\n')
	if err != nil {
		return ans, fmt.Errorf("failed to read from stdin: %w", err)
	}

	return strings.TrimSpace(ans), nil
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup and configure the CDB command line tool.",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL, err := input(
			fmt.Sprintf(
				"What is the base URL of your instance? (default: %s) ",
				config.Current.BaseURL,
			),
		)
		if err != nil {
			return err
		}

		if baseURL != "" {
			config.Current.BaseURL = baseURL
		}

		token, err := login(cmd.Context())
		if err != nil {
			return err
		}

		if token != "" {
			config.Current.Token = token
		}

		filename := config.ConfigFile()
		fh, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer fh.Close()

		err = yaml.NewEncoder(fh).Encode(config.Current)
		if err != nil {
			return err
		}

		fmt.Printf("Wrote new config to %s\n", filename)
		return nil
	},
}
