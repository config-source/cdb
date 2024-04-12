package cdb

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

var (
	ErrConfigValueNotFound = errors.New("config value not found")
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
	Inherited bool      `db:"-"`
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
			return "UNKNOWN STR VALUE"
		}

		return *cv.StrValue
	case TypeInteger:
		if cv.IntValue == nil {
			return math.MaxInt
		}

		return *cv.IntValue
	case TypeFloat:
		if cv.FloatValue == nil {
			return math.MaxFloat32
		}

		return *cv.FloatValue
	case TypeBoolean:
		if cv.BoolValue == nil {
			return false
		}

		return *cv.BoolValue
	default:
		return nil
	}
}

func (cv ConfigValue) String() string {
	return fmt.Sprintf(
		"ConfigValue(id=%d, environment=%d, keyID=%d, name=%s, valueType=%s, value=%v)",
		cv.ID,
		cv.EnvironmentID,
		cv.ConfigKeyID,
		cv.Name,
		cv.ValueType,
		cv.Value(),
	)
}

func (cv ConfigValue) ValueAsString() string {
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
		return ""
	}
}

type ConfigValueRepository interface {
	CreateConfigValue(context.Context, ConfigValue) (ConfigValue, error)

	GetConfiguration(ctx context.Context, environmentName string) ([]ConfigValue, error)
	GetConfigurationValue(ctx context.Context, environmentName, key string) (ConfigValue, error)
}
