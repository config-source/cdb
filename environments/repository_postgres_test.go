package environments_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/config-source/cdb/environments"
	"github.com/config-source/cdb/postgresutils"
	"github.com/rs/zerolog"
)

func initTestDB(t *testing.T) (*environments.Repository, *postgresutils.TestRepository) {
	t.Helper()

	tr, pool := postgresutils.InitTestDB(t)
	repo := environments.NewRepository(
		zerolog.New(nil).Level(zerolog.Disabled),
		pool,
	)

	return repo, tr
}

func envFixture(t *testing.T, repo *environments.Repository, name string, promotesToID *int) environments.Environment {
	env, err := repo.CreateEnvironment(context.Background(), environments.Environment{
		Name:         name,
		PromotesToID: promotesToID,
	})
	if err != nil {
		t.Fatal(err)
	}

	return env
}

func TestCreateEnvironment(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	env, err := repo.CreateEnvironment(context.Background(), environments.Environment{
		Name: "mat",
	})
	if err != nil {
		t.Fatal(err)
	}

	if env.ID == 0 {
		t.Fatalf("Expected ID to be set got: %d", env.ID)
	}

	if env.Name != "mat" {
		t.Fatalf("Expected Name to be mat got: %s", env.Name)
	}

	if env.PromotesToID != nil {
		t.Fatalf("Expected PromotesToID to be nil got: %v", env.PromotesToID)
	}
}

func TestGetEnvironment(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	envFixture(t, repo, "env1", nil)
	env2 := envFixture(t, repo, "env2", nil)

	env, err := repo.GetEnvironment(context.Background(), env2.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(env2, env) {
		t.Fatalf("Got wrong environment expected %v got %v", env2, env)
	}
}

func TestGetEnvironmentReturnsErrEnvNotFound(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	_, err := repo.GetEnvironment(context.Background(), 1)
	if err == nil {
		t.Fatal("Expected an error but got none!")
	}

	if err != environments.ErrNotFound {
		t.Fatalf("Expected an ErrEnvNotFound got: %s", err)
	}
}

func TestGetEnvironmentByName(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	env1 := envFixture(t, repo, "env1", nil)
	envFixture(t, repo, "env2", nil)

	env, err := repo.GetEnvironmentByName(context.Background(), env1.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(env1, env) {
		t.Fatalf("Got wrong environment expected %v got %v", env1, env)
	}
}

func TestGetEnvironmentByNameReturnsErrEnvNotFound(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	_, err := repo.GetEnvironmentByName(context.Background(), "dev")
	if err == nil {
		t.Fatal("Expected an error but got none!")
	}

	if err != environments.ErrNotFound {
		t.Fatalf("Expected an ErrEnvNotFound got: %s", err)
	}
}
