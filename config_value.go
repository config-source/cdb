package cdb

import uuid "github.com/gofrs/uuid/v5"

type ConfigValue struct {
	ID            uuid.UUID `db:"id"`
	ConfigKeyID   uuid.UUID `db:"config_key_id"`
	EnvironmentID uuid.UUID `db:"environment_id"`

	// From config_key tables
	Name      string    `db:"name"`
	ValueType ValueType `db:"value_type"`

	StrValue   *string  `db:"str_value"`
	IntValue   *int     `db:"int_value"`
	FloatValue *float64 `db:"float_value"`
	BoolValue  *bool    `db:"bool_value"`
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
	GetConfiguration(environmentID uuid.UUID) ([]ConfigValue, error)
	GetConfigurationValue(environmentID uuid.UUID, key string) (ConfigValue, error)

	CreateConfigValue(ConfigValue) (ConfigValue, error)
}
