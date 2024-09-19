package configvalues_test

import (
	"context"
	"errors"
	"testing"

	"github.com/config-source/cdb/pkg/auth"
	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/configvalues"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/services"
	"github.com/config-source/cdb/pkg/testutils"
)

func TestServiceCreatesConfigKeyWhenDynamicConfigKeysIsTrue(t *testing.T) {
	promotesToID := 1
	repo := &testutils.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "test",
			},
		},
		Environments: map[int]environments.Environment{
			1: {
				ID:        1,
				Name:      "production",
				ServiceID: 1,
				Service:   "test",
			},
			2: {
				ID:           2,
				Name:         "staging",
				PromotesToID: &promotesToID,
				ServiceID:    1,
				Service:      "test",
			},
		},
		ConfigKeys: map[int]configkeys.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: configkeys.TypeString,
				ServiceID: 1,
				Service:   "test",
			},
			2: {
				ID:        2,
				Name:      "maxReplicas",
				ValueType: configkeys.TypeInteger,
				ServiceID: 1,
				Service:   "test",
			},
		},
	}

	gateway := auth.NewTestGateway()
	service := configvalues.NewService(repo, repo, repo, gateway, true)
	val := 10
	cv, err := service.SetConfigurationValue(
		context.Background(),
		auth.User{},
		2,
		"minReplicas",
		&configvalues.ConfigValue{
			ValueType: configkeys.TypeInteger,
			IntValue:  &val,
		},
	)
	if err != nil {
		t.Fatalf("Failed to set configuration value: %s", err)
	}

	newKey, err := repo.GetConfigKeyByName(context.Background(), 1, "minReplicas")
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
	repo := &testutils.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "test",
			},
		},
		Environments: map[int]environments.Environment{
			1: {
				ID:        1,
				Name:      "production",
				ServiceID: 1,
				Service:   "test",
			},
			2: {
				ID:           2,
				Name:         "staging",
				PromotesToID: &promotesToID,
				ServiceID:    1,
				Service:      "test",
			},
		},
		ConfigKeys: map[int]configkeys.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: configkeys.TypeString,
				ServiceID: 1,
				Service:   "test",
			},
			2: {
				ID:        2,
				Name:      "maxReplicas",
				ValueType: configkeys.TypeInteger,
				ServiceID: 1,
				Service:   "test",
			},
		},
	}

	service := configvalues.NewService(repo, repo, repo, auth.NewTestGateway(), false)
	val := 10
	_, err := service.SetConfigurationValue(
		context.Background(),
		auth.User{},
		2,
		"minReplicas",
		&configvalues.ConfigValue{
			ValueType: configkeys.TypeInteger,
			IntValue:  &val,
		},
	)
	if !errors.Is(err, configkeys.ErrNotFound) {
		t.Fatalf("Expected %s got: %s", configkeys.ErrNotFound, err)
	}
}

func TestServiceReturnsErrorWhenValueTypeIsNotValid(t *testing.T) {
	promotesToID := 1
	repo := &testutils.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "test",
			},
		},
		Environments: map[int]environments.Environment{
			1: {
				ID:        1,
				Name:      "production",
				ServiceID: 1,
				Service:   "test",
			},
			2: {
				ID:           2,
				Name:         "staging",
				PromotesToID: &promotesToID,
				ServiceID:    1,
				Service:      "test",
			},
		},
		ConfigKeys: map[int]configkeys.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: configkeys.TypeString,
				ServiceID: 1,
				Service:   "test",
			},
			2: {
				ID:        2,
				Name:      "maxReplicas",
				ValueType: configkeys.TypeInteger,
				ServiceID: 1,
				Service:   "test",
			},
		},
	}

	service := configvalues.NewService(repo, repo, repo, auth.NewTestGateway(), false)
	val := "test"
	_, err := service.SetConfigurationValue(
		context.Background(),
		auth.User{},
		2,
		"maxReplicas",
		&configvalues.ConfigValue{
			ValueType: configkeys.TypeString,
			StrValue:  &val,
		},
	)
	if !errors.Is(err, configvalues.ErrNotValid) {
		t.Fatalf("Expected %s got: %s", configvalues.ErrNotValid, err)
	}
}
