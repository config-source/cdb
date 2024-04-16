package cdb_test

import (
	"errors"
	"testing"

	"github.com/config-source/cdb"
)

func setBoolValue(cv *cdb.ConfigValue, vt cdb.ValueType) *cdb.ConfigValue {
	cv.SetBoolValue(false)
	cv.ValueType = vt
	return cv

}

func setFloatValue(cv *cdb.ConfigValue, vt cdb.ValueType) *cdb.ConfigValue {
	cv.SetFloatValue(10.10)
	cv.ValueType = vt
	return cv
}

func setIntValue(cv *cdb.ConfigValue, vt cdb.ValueType) *cdb.ConfigValue {
	cv.SetIntValue(10)
	cv.ValueType = vt
	return cv
}

func setStrValue(cv *cdb.ConfigValue, vt cdb.ValueType) *cdb.ConfigValue {
	cv.SetStrValue("test")
	cv.ValueType = vt
	return cv
}

type mutatorFunc func(*cdb.ConfigValue, cdb.ValueType) *cdb.ConfigValue

func TestConfigValueValidatesStrValue(t *testing.T) {
	mutations := []mutatorFunc{
		setBoolValue,
		setIntValue,
		setFloatValue,
	}

	value := cdb.NewStringConfigValue(1, 1, "test")
	if err := value.Valid(); err != nil {
		t.Fatalf("Expected no error got: %s", err)
	}

	for _, mutation := range mutations {
		cv := mutation(value, cdb.TypeString)
		err := cv.Valid()
		if err == nil {
			t.Fatalf("Expected an error for %s got: %s", cv, err)
		}

		if !errors.Is(err, cdb.ErrConfigValueNotValid) {
			t.Fatalf("Expected a cdb.ErrConfigValueNotValid got: %s", err)
		}
	}
}

func TestConfigValueValidatesBoolValue(t *testing.T) {
	mutations := []mutatorFunc{
		setStrValue,
		setIntValue,
		setFloatValue,
	}

	value := cdb.NewBoolConfigValue(1, 1, true)
	if err := value.Valid(); err != nil {
		t.Fatalf("Expected no error got: %s", err)
	}

	for _, mutation := range mutations {
		mutation(value, cdb.TypeBoolean)
		err := value.Valid()
		if err == nil {
			t.Fatalf("Expected an error for %s got: %s", value, err)
		}

		if !errors.Is(err, cdb.ErrConfigValueNotValid) {
			t.Fatalf("Expected a cdb.ErrConfigValueNotValid got: %s", err)
		}
	}
}

func TestConfigValueValidatesIntValue(t *testing.T) {
	mutations := []mutatorFunc{
		setStrValue,
		setBoolValue,
		setFloatValue,
	}

	value := cdb.NewIntConfigValue(1, 1, 10)
	if err := value.Valid(); err != nil {
		t.Fatalf("Expected no error got: %s", err)
	}

	for _, mutation := range mutations {
		cv := mutation(value, cdb.TypeInteger)
		err := cv.Valid()
		if err == nil {
			t.Fatalf("Expected an error for %s got: %s", cv, err)
		}

		if !errors.Is(err, cdb.ErrConfigValueNotValid) {
			t.Fatalf("Expected a cdb.ErrConfigValueNotValid got: %s", err)
		}
	}
}

func TestConfigValueValidatesFloatValue(t *testing.T) {
	mutations := []mutatorFunc{
		setStrValue,
		setBoolValue,
		setIntValue,
	}

	value := cdb.NewFloatConfigValue(1, 1, 10.10)
	if err := value.Valid(); err != nil {
		t.Fatalf("Expected no error got: %s", err)
	}

	for _, mutation := range mutations {
		cv := mutation(value, cdb.TypeFloat)
		err := cv.Valid()
		if err == nil {
			t.Fatalf("Expected an error for %s got: %s", cv, err)
		}

		if !errors.Is(err, cdb.ErrConfigValueNotValid) {
			t.Fatalf("Expected a cdb.ErrConfigValueNotValid got: %s", err)
		}
	}
}
