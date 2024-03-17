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
	Name      string    `db:"name,omitempty"`
	ValueType ValueType `db:"value_type,omitempty"`

	StrValue   *string  `db:"str_value"`
	IntValue   *int     `db:"int_value"`
	FloatValue *float64 `db:"float_value"`
	BoolValue  *bool    `db:"bool_value"`

	CreatedAt time.Time `db:"created_at"`
}

func NewConfigValue(
	valueType ValueType,
	environmentID int,
	configKeyID int,
	strValue *string,
	intValue *int,
	floatValue *float64,
	boolValue *bool,
) ConfigValue {
	return ConfigValue{
		EnvironmentID: environmentID,
		ConfigKeyID:   configKeyID,
		StrValue:      strValue,
		IntValue:      intValue,
		FloatValue:    floatValue,
		BoolValue:     boolValue,
	}
}

func NewBoolConfigValue(environmentID int, configKeyID int, value bool) ConfigValue {
	storedValue := value
	return NewConfigValue(
		TypeBoolean,
		environmentID,
		configKeyID,
		nil,
		nil,
		nil,
		&storedValue,
	)
}

func NewFloatConfigValue(environmentID int, configKeyID int, value float64) ConfigValue {
	storedValue := value
	return NewConfigValue(
		TypeFloat,
		environmentID,
		configKeyID,
		nil,
		nil,
		&storedValue,
		nil,
	)
}

func NewStringConfigValue(environmentID int, configKeyID int, value string) ConfigValue {
	storedValue := value
	return NewConfigValue(
		TypeString,
		environmentID,
		configKeyID,
		&storedValue,
		nil,
		nil,
		nil,
	)
}

func NewIntConfigValue(environmentID int, configKeyID int, value int) ConfigValue {
	storedValue := value
	return NewConfigValue(
		TypeInteger,
		environmentID,
		configKeyID,
		nil,
		&storedValue,
		nil,
		nil,
	)
}

func (cv *ConfigValue) Value() interface{} {
	switch cv.ValueType {
	case TypeString:
		return *cv.StrValue
	case TypeInteger:
		return *cv.IntValue
	case TypeFloat:
		return *cv.FloatValue
	case TypeBoolean:
		return *cv.BoolValue
	default:
		return nil
	}
}

type ConfigValueRepository interface {
	GetConfiguration(ctx context.Context, environmentID int) ([]ConfigValue, error)
	GetConfigurationValue(ctx context.Context, environmentID int, key string) (ConfigValue, error)

	CreateConfigValue(context.Context, ConfigValue) (ConfigValue, error)
}
