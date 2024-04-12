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

func NewConfigValue(environmentID, configKeyID int) *ConfigValue {
	return &ConfigValue{
		EnvironmentID: environmentID,
		ConfigKeyID:   configKeyID,
	}
}

func NewBoolConfigValue(environmentID int, configKeyID int, value bool) ConfigValue {
	return *NewConfigValue(environmentID, configKeyID).SetBoolValue(value)
}

func NewFloatConfigValue(environmentID int, configKeyID int, value float64) ConfigValue {
	return *NewConfigValue(environmentID, configKeyID).SetFloatValue(value)
}

func NewStringConfigValue(environmentID int, configKeyID int, value string) ConfigValue {
	return *NewConfigValue(environmentID, configKeyID).SetStrValue(value)
}

func NewIntConfigValue(environmentID int, configKeyID int, value int) ConfigValue {
	return *NewConfigValue(environmentID, configKeyID).SetIntValue(value)
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

func (cv *ConfigValue) resetValues() *ConfigValue {
	cv.StrValue = nil
	cv.IntValue = nil
	cv.FloatValue = nil
	cv.BoolValue = nil
	return cv
}

func (cv *ConfigValue) SetStrValue(val string) *ConfigValue {
	cv.resetValues()
	storage := val
	cv.StrValue = &storage
	cv.ValueType = TypeString
	return cv
}

func (cv *ConfigValue) SetIntValue(val int) *ConfigValue {
	cv.resetValues()
	storage := val
	cv.IntValue = &storage
	cv.ValueType = TypeInteger
	return cv
}

func (cv *ConfigValue) SetFloatValue(val float64) *ConfigValue {
	cv.resetValues()
	storage := val
	cv.FloatValue = &storage
	cv.ValueType = TypeFloat
	return cv
}

func (cv *ConfigValue) SetBoolValue(val bool) *ConfigValue {
	cv.resetValues()
	storage := val
	cv.BoolValue = &storage
	cv.ValueType = TypeBoolean
	return cv
}

type ConfigValueRepository interface {
	CreateConfigValue(context.Context, ConfigValue) (ConfigValue, error)

	GetConfiguration(ctx context.Context, environmentName string) ([]ConfigValue, error)
	GetConfigurationValue(ctx context.Context, environmentName, key string) (ConfigValue, error)

	UpdateConfigurationValue(context.Context, ConfigValue) (ConfigValue, error)
}
