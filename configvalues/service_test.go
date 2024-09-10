package configvalues_test

import (
	"context"
	"errors"
	"testing"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/configvalues"
	"github.com/config-source/cdb/repository"
)

func TestServiceCreatesConfigKeyWhenDynamicConfigKeysIsTrue(t *testing.T) {
	promotesToID := 1
	repo := &repository.TestRepository{
		Environments: map[int]cdb.Environment{
			1: {
				ID:   1,
				Name: "production",
			},
			2: {
				ID:           2,
				Name:         "staging",
				PromotesToID: &promotesToID,
			},
		},
		ConfigKeys: map[int]cdb.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: cdb.TypeString,
			},
			2: {
				ID:        2,
				Name:      "maxReplicas",
				ValueType: cdb.TypeInteger,
			},
		},
	}

	gateway := auth.NewTestGateway()
	service := configvalues.NewService(repo, gateway, true)
	val := 10
	cv, err := service.SetConfigurationValue(
		context.Background(),
		auth.User{},
		"staging",
		"minReplicas",
		&cdb.ConfigValue{
			ValueType: cdb.TypeInteger,
			IntValue:  &val,
		},
	)
	if err != nil {
		t.Fatalf("Failed to set configuration value: %s", err)
	}

	newKey, err := repo.GetConfigKeyByName(context.Background(), "minReplicas")
	if err != nil {
		t.Fatal(err)
	}

	if cv.ConfigKeyID != newKey.ID {
		t.Fatalf("Expected config value to have the same key ID as the new key: %s %s", newKey, cv)
	}

	if cv.IntValue == nil {
		t.Fatalf("Expected non-nil IntValue: %s", cv)
	}
}

func TestServiceReturnsErrorWhenDynamicConfigKeysIsFalse(t *testing.T) {
	promotesToID := 1
	repo := &repository.TestRepository{
		Environments: map[int]cdb.Environment{
			1: {
				ID:   1,
				Name: "production",
			},
			2: {
				ID:           2,
				Name:         "staging",
				PromotesToID: &promotesToID,
			},
		},
		ConfigKeys: map[int]cdb.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: cdb.TypeString,
			},
			2: {
				ID:        2,
				Name:      "maxReplicas",
				ValueType: cdb.TypeInteger,
			},
		},
	}

	service := configvalues.NewService(repo, auth.NewTestGateway(), false)
	val := 10
	_, err := service.SetConfigurationValue(
		context.Background(),
		auth.User{},
		"staging",
		"minReplicas",
		&cdb.ConfigValue{
			ValueType: cdb.TypeInteger,
			IntValue:  &val,
		},
	)
	if !errors.Is(err, cdb.ErrConfigKeyNotFound) {
		t.Fatalf("Expected %s got: %s", cdb.ErrConfigKeyNotFound, err)
	}
}

func TestServiceReturnsErrorWhenValueTypeIsNotValid(t *testing.T) {
	promotesToID := 1
	repo := &repository.TestRepository{
		Environments: map[int]cdb.Environment{
			1: {
				ID:   1,
				Name: "production",
			},
			2: {
				ID:           2,
				Name:         "staging",
				PromotesToID: &promotesToID,
			},
		},
		ConfigKeys: map[int]cdb.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: cdb.TypeString,
			},
			2: {
				ID:        2,
				Name:      "maxReplicas",
				ValueType: cdb.TypeInteger,
			},
		},
	}

	service := configvalues.NewService(repo, auth.NewTestGateway(), false)
	val := "test"
	_, err := service.SetConfigurationValue(
		context.Background(),
		auth.User{},
		"staging",
		"maxReplicas",
		&cdb.ConfigValue{
			ValueType: cdb.TypeString,
			StrValue:  &val,
		},
	)
	if !errors.Is(err, cdb.ErrConfigValueNotValid) {
		t.Fatalf("Expected %s got: %s", cdb.ErrConfigValueNotValid, err)
	}
}
