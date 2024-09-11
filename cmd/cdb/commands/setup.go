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
		var err error
		if config.Current.BaseURL == "" {
			config.Current.BaseURL, err = input("What is the base URL of your instance? ")
		} else {
			var response string
			response, err = input(
				fmt.Sprintf(
					"What is the base URL of your instance? (default: %s) ",
					config.Current.BaseURL,
				),
			)

			if response != "" {
				config.Current.BaseURL = response
			}
		}

		if err != nil {
			return err
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

		// TODO: generate a token or tell a user how to do it.
		return nil

	},
}
