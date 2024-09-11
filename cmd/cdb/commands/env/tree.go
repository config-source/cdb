package env

import (
	"context"
	"fmt"
	"strings"

	"github.com/config-source/cdb/cmd/cdb/config"
	"github.com/config-source/cdb/environments"
	"github.com/spf13/cobra"
)

func printTree(tree environments.Tree, depth int) {
	indent := ""
	leading := strings.Repeat("─", depth)
	parentMarker := ""
	if depth != 0 {
		parentMarker = "└"
		indent = strings.Repeat("   ", depth)
	}

	fmt.Printf("%s%s%s %s\n", indent, parentMarker, leading, tree.Environment.Name)

	for _, child := range tree.Children {
		printTree(child, depth+1)
	}
}

var envTreeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Print the promotion tree of your environments",
	RunE: func(cmd *cobra.Command, args []string) error {
		trees, err := config.Client.GetEnvironmentTree(context.Background())
		if err != nil {
			return err
		}

		for _, tree := range trees {
			printTree(tree, 0)
		}

		return nil
	},
}

func init() {
	Command.AddCommand(envTreeCmd)
}
