package configkeys_test

import (
	"context"
	_ "embed"
	"reflect"
	"testing"

	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/postgresutils"
	"github.com/config-source/cdb/pkg/services"
	"github.com/rs/zerolog"
)

func initTestDB(t *testing.T) (configkeys.Repository, services.Repository, *postgresutils.TestDatabase) {
	t.Helper()

	tr, pool := postgresutils.InitTestDB(t)

	repo := configkeys.NewRepository(
		zerolog.New(nil).Level(zerolog.Disabled),
		pool,
	)
	svcRepo := services.NewRepository(
		zerolog.New(nil).Level(zerolog.Disabled),
		pool,
	)

	return repo, svcRepo, tr
}

func svcFixture(t *testing.T, repo services.Repository, name string) services.Service {
	svc, err := repo.CreateService(context.Background(), services.Service{
		Name: name,
	})
	if err != nil {
		t.Fatal(err)
	}

	return svc
}

func configKeyFixture(
	t *testing.T,
	repo configkeys.Repository,
	svcID int,
	name string,
	valueType configkeys.ValueType,
	canPropagate bool,
) configkeys.ConfigKey {
	ck, err := repo.CreateConfigKey(
		context.Background(),
		configkeys.NewWithPropagation(
			svcID,
			name,
			valueType,
			canPropagate,
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	return ck
}

func TestCreateConfigKey(t *testing.T) {
	repo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "test")
	ck, err := repo.CreateConfigKey(context.Background(), configkeys.New(
		svc.ID,
		"mat",
		configkeys.TypeString,
	))
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

	ck2, err := repo.CreateConfigKey(
		context.Background(),
		configkeys.NewWithPropagation(
			svc.ID,
			"mat2",
			configkeys.TypeString,
			false,
		),
	)
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
	repo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "test")
	configKey1 := configKeyFixture(
		t,
		repo,
		svc.ID,
		"getConfigKey1",
		configkeys.TypeInteger,
		true,
	)
	configKey1.Service = svc.Name
	configKeyFixture(
		t,
		repo,
		svc.ID,
		"getConfigKey2",
		configkeys.TypeString,
		true,
	)

	configKey, err := repo.GetConfigKey(context.Background(), configKey1.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(configKey1, configKey) {
		t.Fatalf("Got wrong configKey expected %v got %v", configKey1, configKey)
	}
}

func TestListConfigKeys(t *testing.T) {
	repo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "test")
	svc2 := svcFixture(t, svcRepo, "test2")

	ck1 := configKeyFixture(t, repo, svc.ID, "configKey1", configkeys.TypeInteger, true)
	ck2 := configKeyFixture(t, repo, svc.ID, "configKey2", configkeys.TypeString, true)
	ck3 := configKeyFixture(t, repo, svc2.ID, "configKey3", configkeys.TypeBoolean, true)

	ck1.Service = svc.Name
	ck2.Service = svc.Name
	ck3.Service = svc.Name

	configKeys := []configkeys.ConfigKey{
		ck1,
		ck2,
	}

	retrieved, err := repo.ListConfigKeys(context.Background(), svc.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(retrieved, configKeys) {
		t.Fatalf("Expected config keys: %v Got: %v", configKeys, retrieved)
	}
}

func TestListConfigKeysNoService(t *testing.T) {
	repo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc1 := svcFixture(t, svcRepo, "test1")
	svc2 := svcFixture(t, svcRepo, "test2")
	svc3 := svcFixture(t, svcRepo, "test3")

	ck1 := configKeyFixture(t, repo, svc1.ID, "configKey1", configkeys.TypeInteger, true)
	ck2 := configKeyFixture(t, repo, svc2.ID, "configKey2", configkeys.TypeString, true)
	ck3 := configKeyFixture(t, repo, svc3.ID, "configKey3", configkeys.TypeBoolean, true)

	ck1.Service = svc1.Name
	ck2.Service = svc2.Name
	ck3.Service = svc3.Name

	configKeys := []configkeys.ConfigKey{
		ck1,
		ck2,
		ck3,
	}

	retrieved, err := repo.ListConfigKeys(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(retrieved, configKeys) {
		t.Fatalf("Expected config keys: %v Got: %v", configKeys, retrieved)
	}
}
