package cdb

import (
	"context"
	"fmt"
	"time"
)

type ConfigValue struct {
	ID            int `db:"id"`
	ConfigKeyID   int `db:"config_key_id"`
	EnvironmentID int `db:"environment_id"`

	// From config_key tables
	Name      string    `db:"name"`
	ValueType ValueType `db:"value_type"`

	StrValue   *string  `db:"str_value"`
	IntValue   *int     `db:"int_value"`
	FloatValue *float64 `db:"float_value"`
	BoolValue  *bool    `db:"bool_value"`

	CreatedAt time.Time `db:"created_at"`
}

func NewBoolConfigValue(environmentID int, configKeyID int, value bool) ConfigValue {
	storedValue := value
	return ConfigValue{
		ValueType:     TypeBoolean,
		EnvironmentID: environmentID,
		ConfigKeyID:   configKeyID,
		BoolValue:     &storedValue,
	}
}

func NewFloatConfigValue(environmentID int, configKeyID int, value float64) ConfigValue {
	storedValue := value
	return ConfigValue{
		ValueType:     TypeFloat,
		EnvironmentID: environmentID,
		ConfigKeyID:   configKeyID,
		FloatValue:    &storedValue,
	}
}

func NewStringConfigValue(environmentID int, configKeyID int, value string) ConfigValue {
	storedValue := value
	return ConfigValue{
		ValueType:     TypeFloat,
		EnvironmentID: environmentID,
		ConfigKeyID:   configKeyID,
		StrValue:      &storedValue,
	}
}

func NewIntConfigValue(environmentID int, configKeyID int, value int) ConfigValue {
	storedValue := value
	return ConfigValue{
		ValueType:     TypeFloat,
		EnvironmentID: environmentID,
		ConfigKeyID:   configKeyID,
		IntValue:      &storedValue,
	}
}

func (cv *ConfigValue) Value() interface{} {
	switch cv.ValueType {
	case TypeString:
		if cv.StrValue == nil {
			panic("UNKNOWN STR VALUE")
		}

		return *cv.StrValue
	case TypeInteger:
		if cv.IntValue == nil {
			panic("UNKNOWN INT VALUE")
		}

		return *cv.IntValue
	case TypeFloat:
		if cv.FloatValue == nil {
			panic("UNKNOWN FLOAT VALUE")
		}

		return *cv.FloatValue
	case TypeBoolean:
		if cv.BoolValue == nil {
			panic("UNKNOWN BOOLEAN VALUE")
		}

		return *cv.BoolValue
	default:
		return nil
	}
}

func (cv ConfigValue) String() string {
	return fmt.Sprintf(
		"ConfigValue(%d, %d, %s, %v)",
		cv.EnvironmentID,
		cv.ConfigKeyID,
		cv.Name,
		cv.Value(),
	)
}

type ConfigValueRepository interface {
	CreateConfigValue(context.Context, ConfigValue) (ConfigValue, error)

	GetConfiguration(ctx context.Context, environmentID int) ([]ConfigValue, error)
	GetConfigurationValue(ctx context.Context, environmentName, key string) (ConfigValue, error)
}
