package cdb_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/config-source/cdb"
)

func setBoolValue(cv *cdb.ConfigValue) *cdb.ConfigValue {
	return cv.SetBoolValue(false)
}

func setFloatValue(cv *cdb.ConfigValue) *cdb.ConfigValue {
	return cv.SetFloatValue(10.10)
}

func setIntValue(cv *cdb.ConfigValue) *cdb.ConfigValue {
	return cv.SetIntValue(10)
}

func setStrValue(cv *cdb.ConfigValue) *cdb.ConfigValue {
	return cv.SetStrValue("test")
}

type mutatorFunc func(*cdb.ConfigValue) *cdb.ConfigValue

func TestConfigValueValidatesStrValue(t *testing.T) {
	reset := setStrValue
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
		cv := mutation(reset(value))
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
	// reset := setBoolValue
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
		cv := mutation(value)
		fmt.Println("outside", cv)
		err := cv.Valid()
		if err == nil {
			t.Fatalf("Expected an error for %s got: %s", cv, err)
		}

		if !errors.Is(err, cdb.ErrConfigValueNotValid) {
			t.Fatalf("Expected a cdb.ErrConfigValueNotValid got: %s", err)
		}
	}
}

func TestConfigValueValidatesIntValue(t *testing.T) {
	reset := setIntValue
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
		cv := mutation(reset(value))
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
	reset := setFloatValue
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
		cv := mutation(reset(value))
		err := cv.Valid()
		if err == nil {
			t.Fatalf("Expected an error for %s got: %s", cv, err)
		}

		if !errors.Is(err, cdb.ErrConfigValueNotValid) {
			t.Fatalf("Expected a cdb.ErrConfigValueNotValid got: %s", err)
		}
	}
}
