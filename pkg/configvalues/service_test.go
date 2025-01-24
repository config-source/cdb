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
)

func setupBasicService(t *testing.T, tc TestContext) {
	t.Helper()

	svc, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}

	production, err := tc.environmentRepo.CreateEnvironment(context.Background(), environments.Environment{Name: "production", ServiceID: svc.ID})
	if err != nil {
		t.Fatal(err)
	}

	_, err = tc.environmentRepo.CreateEnvironment(
		context.Background(),
		environments.Environment{
			Name:         "staging",
			ServiceID:    svc.ID,
			PromotesToID: &production.ID,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, key := range []configkeys.ConfigKey{
		{
			Name:      "owner",
			ValueType: configkeys.TypeString,
			ServiceID: svc.ID,
		},
		{
			Name:      "maxReplicas",
			ValueType: configkeys.TypeInteger,
			ServiceID: svc.ID,
		},
	} {
		_, err := tc.keyRepo.CreateConfigKey(context.Background(), key)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestServiceCreatesConfigKeyWhenDynamicConfigKeysIsTrue(t *testing.T) {
	tc := initTestDB(t)
	setupBasicService(t, tc)

	gateway := auth.NewTestGateway()
	service := configvalues.NewService(tc.valueRepo, tc.environmentRepo, tc.keyRepo, gateway, true)
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

	newKey, err := tc.keyRepo.GetConfigKeyByName(context.Background(), "test", "minReplicas")
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
	tc := initTestDB(t)
	setupBasicService(t, tc)

	service := configvalues.NewService(
		tc.valueRepo,
		tc.environmentRepo,
		tc.keyRepo,
		auth.NewTestGateway(),
		false,
	)
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
	tc := initTestDB(t)
	setupBasicService(t, tc)

	service := configvalues.NewService(
		tc.valueRepo,
		tc.environmentRepo,
		tc.keyRepo,
		auth.NewTestGateway(),
		false,
	)
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
