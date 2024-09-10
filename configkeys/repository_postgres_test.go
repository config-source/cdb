package configkeys_test

import (
	"context"
	_ "embed"
	"reflect"
	"testing"

	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/postgresutils"
	"github.com/rs/zerolog"
)

func initTestDB(t *testing.T) (*configkeys.PostgresRepository, *postgresutils.TestRepository) {
	t.Helper()

	tr, pool := postgresutils.InitTestDB(t)

	repo := configkeys.NewRepository(
		zerolog.New(nil).Level(zerolog.Disabled),
		pool,
	)

	return repo, tr
}

func configKeyFixture(t *testing.T, repo *configkeys.PostgresRepository, name string, valueType configkeys.ValueType, canPropagate bool) configkeys.ConfigKey {
	ck, err := repo.CreateConfigKey(context.Background(), configkeys.ConfigKey{
		Name:         name,
		ValueType:    valueType,
		CanPropagate: &canPropagate,
	})
	if err != nil {
		t.Fatal(err)
	}

	return ck
}

func TestCreateConfigKey(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	ck, err := repo.CreateConfigKey(context.Background(), configkeys.ConfigKey{
		Name:      "mat",
		ValueType: configkeys.TypeString,
	})
	if err != nil {
		t.Fatal(err)
	}

	if ck.ID == 0 {
		t.Fatalf("Expected ID to be set got: %d", ck.ID)
	}

	if ck.Name != "mat" {
		t.Fatalf("Expected Name to be mat got: %s", ck.Name)
	}

	if ck.ValueType != configkeys.TypeString {
		t.Fatalf("Expected ValueType to be %d got: %d", configkeys.TypeString, ck.ValueType)
	}

	if !*ck.CanPropagate {
		t.Fatalf("Expected can propagate to default to true got: %v", *ck.CanPropagate)
	}

	ck2, err := repo.CreateConfigKey(context.Background(), configkeys.NewWithPropagation("mat2", configkeys.TypeString, false))
	if err != nil {
		t.Fatal(err)
	}

	if ck2.ID == 0 {
		t.Fatalf("Expected ID to be set got: %d", ck2.ID)
	}

	if ck2.Name != "mat2" {
		t.Fatalf("Expected Name to be mat got: %s", ck2.Name)
	}

	if ck2.ValueType != configkeys.TypeString {
		t.Fatalf("Expected ValueType to be %d got: %d", configkeys.TypeString, ck2.ValueType)
	}

	if *ck2.CanPropagate {
		t.Fatalf("Expected can propagate to to be false got: %v", *ck2.CanPropagate)
	}
}

func TestGetConfigKey(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	configKey1 := configKeyFixture(t, repo, "getConfigKey1", configkeys.TypeInteger, true)
	configKeyFixture(t, repo, "getConfigKey2", configkeys.TypeString, true)

	configKey, err := repo.GetConfigKey(context.Background(), configKey1.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(configKey1, configKey) {
		t.Fatalf("Got wrong configKeyironment expected %v got %v", configKey1, configKey)
	}
}

func TestListConfigKeys(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	configKeys := []configkeys.ConfigKey{
		configKeyFixture(t, repo, "configKey1", configkeys.TypeInteger, true),
		configKeyFixture(t, repo, "configKey2", configkeys.TypeString, true),
		configKeyFixture(t, repo, "configKey3", configkeys.TypeBoolean, true),
	}

	retrieved, err := repo.ListConfigKeys(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(retrieved, configKeys) {
		t.Fatalf("Expected config keys: %v Got: %v", configKeys, retrieved)
	}
}
