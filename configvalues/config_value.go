package configvalues

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/config-source/cdb/configkeys"
)

var (
	ErrNotFound           = errors.New("config value not found")
	ErrNotValid           = errors.New("config value is not valid")
	ErrAlreadySet         = errors.New("config value is already set for this environment")
	ErrValueTypeMustBeSet = errors.New("must set ValueType on the config value when trying to dynamically create a key")
)

type ConfigValue struct {
	ID            int `db:"id"`
	ConfigKeyID   int `db:"config_key_id"`
	EnvironmentID int `db:"environment_id"`

	// From config_key tables
	Name      string               `db:"name"`
	ValueType configkeys.ValueType `db:"value_type"`

	StrValue   *string  `db:"str_value"`
	IntValue   *int     `db:"int_value"`
	FloatValue *float64 `db:"float_value"`
	BoolValue  *bool    `db:"bool_value"`

	CreatedAt time.Time `db:"created_at"`
	// Inherited indicates that the value was inherited
	Inherited bool `db:"-"`
	// InheritedFrom is the name of the Enviroment that the value was inherited from.
	InheritedFrom string `db:"-"`
}

func NewConfigValue(environmentID, configKeyID int) *ConfigValue {
	return &ConfigValue{
		EnvironmentID: environmentID,
		ConfigKeyID:   configKeyID,
	}
}

func NewBoolConfigValue(environmentID int, configKeyID int, value bool) *ConfigValue {
	return NewConfigValue(environmentID, configKeyID).SetBoolValue(value)
}

func NewFloatConfigValue(environmentID int, configKeyID int, value float64) *ConfigValue {
	return NewConfigValue(environmentID, configKeyID).SetFloatValue(value)
}

func NewStringConfigValue(environmentID int, configKeyID int, value string) *ConfigValue {
	return NewConfigValue(environmentID, configKeyID).SetStrValue(value)
}

func NewIntConfigValue(environmentID int, configKeyID int, value int) *ConfigValue {
	return NewConfigValue(environmentID, configKeyID).SetIntValue(value)
}

func (cv *ConfigValue) Value() interface{} {
	if err := cv.Valid(); err != nil {
		return err
	}

	switch cv.ValueType {
	case configkeys.TypeString:
		return *cv.StrValue
	case configkeys.TypeInteger:
		return *cv.IntValue
	case configkeys.TypeFloat:
		return *cv.FloatValue
	case configkeys.TypeBoolean:
		return *cv.BoolValue
	// This should be unreachable.
	default:
		fmt.Println("ConfigValue somehow reached unreachable code!")
		return nil
	}
}

func (cv *ConfigValue) String() string {
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

func (cv *ConfigValue) ValueAsString() string {
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
	cv.ValueType = configkeys.TypeString
	return cv
}

func (cv *ConfigValue) SetIntValue(val int) *ConfigValue {
	cv.resetValues()
	storage := val
	cv.IntValue = &storage
	cv.ValueType = configkeys.TypeInteger
	return cv
}

func (cv *ConfigValue) SetFloatValue(val float64) *ConfigValue {
	cv.resetValues()
	storage := val
	cv.FloatValue = &storage
	cv.ValueType = configkeys.TypeFloat
	return cv
}

func (cv *ConfigValue) SetBoolValue(val bool) *ConfigValue {
	cv.resetValues()
	storage := val
	cv.BoolValue = &storage
	cv.ValueType = configkeys.TypeBoolean
	return cv
}

func (cv *ConfigValue) Valid() error {
	switch cv.ValueType {
	case configkeys.TypeBoolean:
		return cv.validateBoolean()
	case configkeys.TypeFloat:
		return cv.validateFloat()
	case configkeys.TypeInteger:
		return cv.validateInt()
	case configkeys.TypeString:
		return cv.validateStr()
	default:
		return fmt.Errorf("%w: unrecognised ValueType: %s", ErrNotValid, cv.ValueType)
	}
}

func (cv *ConfigValue) validateStr() error {
	if cv.IntValue != nil {
		return fmt.Errorf("%w: IntValue must be null for string ConfigValue", ErrNotValid)
	}

	if cv.FloatValue != nil {
		return fmt.Errorf("%w: FloatValue must be null for string ConfigValue", ErrNotValid)
	}

	if cv.BoolValue != nil {
		return fmt.Errorf("%w: BoolValue must be null for string ConfigValue", ErrNotValid)
	}

	if cv.StrValue == nil {
		return fmt.Errorf("%w: StrValue must not be null for string ConfigValue", ErrNotValid)
	}

	return nil
}

func (cv *ConfigValue) validateInt() error {
	if cv.StrValue != nil {
		return fmt.Errorf("%w: StrValue must be null for int ConfigValue", ErrNotValid)
	}

	if cv.FloatValue != nil {
		return fmt.Errorf("%w: FloatValue must be null for int ConfigValue", ErrNotValid)
	}

	if cv.BoolValue != nil {
		return fmt.Errorf("%w: BoolValue must be null for int ConfigValue", ErrNotValid)
	}

	if cv.IntValue == nil {
		return fmt.Errorf("%w: IntValue must not be null for int ConfigValue", ErrNotValid)
	}

	return nil
}

func (cv *ConfigValue) validateFloat() error {
	if cv.StrValue != nil {
		return fmt.Errorf("%w: StrValue must be null for float ConfigValue", ErrNotValid)
	}

	if cv.IntValue != nil {
		return fmt.Errorf("%w: IntValue must be null for float ConfigValue", ErrNotValid)
	}

	if cv.BoolValue != nil {
		return fmt.Errorf("%w: BoolValue must be null for float ConfigValue", ErrNotValid)
	}

	if cv.FloatValue == nil {
		return fmt.Errorf("%w: FloatValue must not be null for float ConfigValue", ErrNotValid)
	}

	return nil
}

func (cv *ConfigValue) validateBoolean() error {
	if cv.StrValue != nil {
		return fmt.Errorf("%w: StrValue must be null for boolean ConfigValue", ErrNotValid)
	}

	if cv.IntValue != nil {
		return fmt.Errorf("%w: IntValue must be null for boolean ConfigValue", ErrNotValid)
	}

	if cv.FloatValue != nil {
		return fmt.Errorf("%w: FloatValue must be null for boolean ConfigValue", ErrNotValid)
	}

	if cv.BoolValue == nil {
		return fmt.Errorf("%w: BoolValue must not be null for boolean ConfigValue", ErrNotValid)
	}

	return nil
}

type ConfigValueRepository interface {
	CreateConfigValue(context.Context, *ConfigValue) (*ConfigValue, error)

	GetConfiguration(ctx context.Context, environmentName string) ([]ConfigValue, error)
	GetConfigurationValue(ctx context.Context, environmentName, key string) (*ConfigValue, error)
	GetConfigValueByEnvAndKey(ctx context.Context, environmentName, key string) (*ConfigValue, error)

	UpdateConfigurationValue(context.Context, *ConfigValue) (*ConfigValue, error)
}
