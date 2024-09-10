package configvalues_test

import (
	"errors"
	"testing"

	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/configvalues"
)

func setBoolValue(cv *configvalues.ConfigValue, vt configkeys.ValueType) *configvalues.ConfigValue {
	cv.SetBoolValue(false)
	cv.ValueType = vt
	return cv

}

func setFloatValue(cv *configvalues.ConfigValue, vt configkeys.ValueType) *configvalues.ConfigValue {
	cv.SetFloatValue(10.10)
	cv.ValueType = vt
	return cv
}

func setIntValue(cv *configvalues.ConfigValue, vt configkeys.ValueType) *configvalues.ConfigValue {
	cv.SetIntValue(10)
	cv.ValueType = vt
	return cv
}

func setStrValue(cv *configvalues.ConfigValue, vt configkeys.ValueType) *configvalues.ConfigValue {
	cv.SetStrValue("test")
	cv.ValueType = vt
	return cv
}

type mutatorFunc func(*configvalues.ConfigValue, configkeys.ValueType) *configvalues.ConfigValue

func TestConfigValueValidatesStrValue(t *testing.T) {
	mutations := []mutatorFunc{
		setBoolValue,
		setIntValue,
		setFloatValue,
	}

	value := configvalues.NewStringConfigValue(1, 1, "test")
	if err := value.Valid(); err != nil {
		t.Fatalf("Expected no error got: %s", err)
	}

	for _, mutation := range mutations {
		cv := mutation(value, configkeys.TypeString)
		err := cv.Valid()
		if err == nil {
			t.Fatalf("Expected an error for %s got: %s", cv, err)
		}

		if !errors.Is(err, configvalues.ErrNotValid) {
			t.Fatalf("Expected a configvalues.ErrConfigValueNotValid got: %s", err)
		}
	}
}

func TestConfigValueValidatesBoolValue(t *testing.T) {
	mutations := []mutatorFunc{
		setStrValue,
		setIntValue,
		setFloatValue,
	}

	value := configvalues.NewBoolConfigValue(1, 1, true)
	if err := value.Valid(); err != nil {
		t.Fatalf("Expected no error got: %s", err)
	}

	for _, mutation := range mutations {
		mutation(value, configkeys.TypeBoolean)
		err := value.Valid()
		if err == nil {
			t.Fatalf("Expected an error for %s got: %s", value, err)
		}

		if !errors.Is(err, configvalues.ErrNotValid) {
			t.Fatalf("Expected a configvalues.ErrConfigValueNotValid got: %s", err)
		}
	}
}

func TestConfigValueValidatesIntValue(t *testing.T) {
	mutations := []mutatorFunc{
		setStrValue,
		setBoolValue,
		setFloatValue,
	}

	value := configvalues.NewIntConfigValue(1, 1, 10)
	if err := value.Valid(); err != nil {
		t.Fatalf("Expected no error got: %s", err)
	}

	for _, mutation := range mutations {
		cv := mutation(value, configkeys.TypeInteger)
		err := cv.Valid()
		if err == nil {
			t.Fatalf("Expected an error for %s got: %s", cv, err)
		}

		if !errors.Is(err, configvalues.ErrNotValid) {
			t.Fatalf("Expected a configvalues.ErrConfigValueNotValid got: %s", err)
		}
	}
}

func TestConfigValueValidatesFloatValue(t *testing.T) {
	mutations := []mutatorFunc{
		setStrValue,
		setBoolValue,
		setIntValue,
	}

	value := configvalues.NewFloatConfigValue(1, 1, 10.10)
	if err := value.Valid(); err != nil {
		t.Fatalf("Expected no error got: %s", err)
	}

	for _, mutation := range mutations {
		cv := mutation(value, configkeys.TypeFloat)
		err := cv.Valid()
		if err == nil {
			t.Fatalf("Expected an error for %s got: %s", cv, err)
		}

		if !errors.Is(err, configvalues.ErrNotValid) {
			t.Fatalf("Expected a configvalues.ErrConfigValueNotValid got: %s", err)
		}
	}
}
