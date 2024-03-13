package cdb

import (
	"context"
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

func (cv *ConfigValue) Value() interface{} {
	switch cv.ValueType {
	case TypeString:
		return cv.StrValue
	case TypeInteger:
		return cv.IntValue
	case TypeFloat:
		return cv.FloatValue
	case TypeBoolean:
		return cv.BoolValue
	default:
		return nil
	}
}

type ConfigValueRepository interface {
	GetConfiguration(ctx context.Context, environmentID int) ([]ConfigValue, error)
	GetConfigurationValue(ctx context.Context, int, key string) (ConfigValue, error)

	CreateConfigValue(context.Context, ConfigValue) (ConfigValue, error)
}
