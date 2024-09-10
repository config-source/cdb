package postgres_test

import (
	"context"
	_ "embed"
	"reflect"
	"testing"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/repository/postgres"
)

func configKeyFixture(t *testing.T, repo *postgres.Repository, name string, valueType cdb.ValueType, canPropagate bool) cdb.ConfigKey {
	ck, err := repo.CreateConfigKey(context.Background(), cdb.ConfigKey{
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

	ck, err := repo.CreateConfigKey(context.Background(), cdb.ConfigKey{
		Name:      "mat",
		ValueType: cdb.TypeString,
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

	if ck.ValueType != cdb.TypeString {
		t.Fatalf("Expected ValueType to be %d got: %d", cdb.TypeString, ck.ValueType)
	}

	if !*ck.CanPropagate {
		t.Fatalf("Expected can propagate to default to true got: %v", *ck.CanPropagate)
	}

	ck2, err := repo.CreateConfigKey(context.Background(), cdb.NewConfigKeyWithCanPropagate("mat2", cdb.TypeString, false))
	if err != nil {
		t.Fatal(err)
	}

	if ck2.ID == 0 {
		t.Fatalf("Expected ID to be set got: %d", ck2.ID)
	}

	if ck2.Name != "mat2" {
		t.Fatalf("Expected Name to be mat got: %s", ck2.Name)
	}

	if ck2.ValueType != cdb.TypeString {
		t.Fatalf("Expected ValueType to be %d got: %d", cdb.TypeString, ck2.ValueType)
	}

	if *ck2.CanPropagate {
		t.Fatalf("Expected can propagate to to be false got: %v", *ck2.CanPropagate)
	}
}

func TestGetConfigKey(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	configKey1 := configKeyFixture(t, repo, "getConfigKey1", cdb.TypeInteger, true)
	configKeyFixture(t, repo, "getConfigKey2", cdb.TypeString, true)

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

	configKeys := []cdb.ConfigKey{
		configKeyFixture(t, repo, "configKey1", cdb.TypeInteger, true),
		configKeyFixture(t, repo, "configKey2", cdb.TypeString, true),
		configKeyFixture(t, repo, "configKey3", cdb.TypeBoolean, true),
	}

	retrieved, err := repo.ListConfigKeys(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(retrieved, configKeys) {
		t.Fatalf("Expected config keys: %v Got: %v", configKeys, retrieved)
	}
}
