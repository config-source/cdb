package configuration

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"github.com/config-source/cdb/cmd/cdb/config"
	"github.com/config-source/cdb/cmd/cdb/table"
	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/configvalues"
	"github.com/spf13/cobra"
)

func valueToRow(cv configvalues.ConfigValue) []string {
	repr := ""
	switch cv.ValueType {
	case configkeys.TypeString:
		repr = *cv.StrValue
	case configkeys.TypeInteger:
		repr = fmt.Sprintf("%d", *cv.IntValue)
	case configkeys.TypeFloat:
		repr = fmt.Sprintf("%f", *cv.FloatValue)
	case configkeys.TypeBoolean:
		repr = fmt.Sprintf("%t", *cv.BoolValue)
	default:
		repr = "UNKNOWN VALUE!"
	}

	return []string{
		cv.Name,
		repr,
		fmt.Sprintf("%t", cv.Inherited),
	}
}

func printConfigTable(values []configvalues.ConfigValue) {
	tbl := table.Table{
		Headings: []string{"Key", "Value", "Inherited"},
		Rows:     make([][]string, len(values)),
	}

	slices.SortFunc(values, func(a, b configvalues.ConfigValue) int {
		return cmp.Compare(a.Name, b.Name)
	})

	for idx, value := range values {
		tbl.Rows[idx] = valueToRow(value)
	}

	fmt.Println(tbl)
}

var getConfigCmd = &cobra.Command{
	Use: "get <environment-name> [configuration-key-name]",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		var env, key string
		env = args[0]
		if len(args) > 1 {
			key = args[1]
		}

		if key != "" {
			value, err := config.Client.GetConfigurationValue(ctx, env, key)
			if err != nil {
				return err
			}

			fmt.Println(value.Value())
		} else {
			values, err := config.Client.GetConfiguration(ctx, env)
			if err != nil {
				return err
			}

			printConfigTable(values)
		}

		return nil
	},
}
