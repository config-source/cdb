package configuration

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/config-source/cdb/cmd/cdb/config"
	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/configvalues"
	"github.com/spf13/cobra"
)

var (
	env   string
	key   string
	value string
)

func valueAsString(cv *configvalues.ConfigValue) string {
	switch v := cv.Value().(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case float64, float32:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		fmt.Println(cv)
		return ""
	}
}

func determineValueType(value string) (interface{}, configkeys.ValueType) {
	intVal, err := strconv.Atoi(value)
	if err == nil {
		return intVal, configkeys.TypeInteger
	}

	floatVal, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return floatVal, configkeys.TypeFloat
	}

	lowered := strings.ToLower(value)
	if lowered == "true" || lowered == "false" {
		return lowered == "true", configkeys.TypeBoolean
	}

	return value, configkeys.TypeString
}

var setConfigCmd = &cobra.Command{
	Use: "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		configValue := &configvalues.ConfigValue{}
		value, valueType := determineValueType(value)

		switch valueType {
		case configkeys.TypeString:
			configValue.SetStrValue(value.(string))
		case configkeys.TypeBoolean:
			configValue.SetBoolValue(value.(bool))
		case configkeys.TypeInteger:
			configValue.SetIntValue(value.(int))
		case configkeys.TypeFloat:
			configValue.SetFloatValue(value.(float64))
		default:
			return errors.New("somehow couldn't find the data type of the config key")
		}

		if err := configValue.Valid(); err != nil {
			return err
		}

		created, err := config.Client.SetConfigurationValue(context.Background(), env, key, configValue)
		if err != nil {
			return err
		}

		fmt.Printf("Set %s=%s for %s\n", key, valueAsString(created), env)
		return nil
	},
}

func init() {
	setConfigCmd.Flags().StringVarP(&env, "environment", "e", "", "The environment you want to set the value for, accepts an environment name or ID.")
	setConfigCmd.Flags().StringVarP(&key, "key", "k", "", "The configuration key you want to set the value for, accepts a key name or ID.")
	setConfigCmd.Flags().StringVarP(&value, "value", "v", "", "The value you want to set the config key to.")
	setConfigCmd.MarkFlagRequired("environment") // nolint:errcheck
	setConfigCmd.MarkFlagRequired("key")         // nolint:errcheck
	setConfigCmd.MarkFlagRequired("value")       // nolint:errcheck

	Command.AddCommand(setConfigCmd)
}
