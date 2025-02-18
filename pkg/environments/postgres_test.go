package environments_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/postgresutils"
	"github.com/config-source/cdb/pkg/services"
	"github.com/rs/zerolog"
)

func initTestDB(t *testing.T) (*environments.Repository, *services.Repository) {
	t.Helper()

	pool := postgresutils.InitTestDB(t)
	repo := environments.NewRepository(
		zerolog.New(nil).Level(zerolog.Disabled),
		pool,
	)
	svcRepo := services.NewRepository(
		zerolog.New(nil).Level(zerolog.Disabled),
		pool,
	)

	return repo, svcRepo
}

func envFixture(
	t *testing.T,
	repo *environments.Repository,
	name string,
	promotesToID *int,
	serviceID int,
) environments.Environment {
	env, err := repo.CreateEnvironment(context.Background(), environments.Environment{
		Name:         name,
		PromotesToID: promotesToID,
		ServiceID:    serviceID,
	})
	if err != nil {
		t.Fatal(err)
	}

	return env
}

func svcFixture(t *testing.T, repo *services.Repository, name string) services.Service {
	svc, err := repo.CreateService(context.Background(), services.Service{
		Name: name,
	})
	if err != nil {
		t.Fatal(err)
	}

	return svc
}

func TestCreateEnvironment(t *testing.T) {
	repo, svcRepo := initTestDB(t)

	svc := svcFixture(t, svcRepo, "svc1")

	env, err := repo.CreateEnvironment(context.Background(), environments.Environment{
		Name:      "mat",
		ServiceID: svc.ID,
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
	repo, svcRepo := initTestDB(t)

	svc := svcFixture(t, svcRepo, "svc1")

	envFixture(t, repo, "env1", nil, svc.ID)
	env2 := envFixture(t, repo, "env2", nil, svc.ID)
	env2.Service = svc.Name

	env, err := repo.GetEnvironment(context.Background(), env2.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(env2, env) {
		t.Fatalf("Got wrong environment expected %v got %v", env2, env)
	}
}

func TestGetEnvironmentReturnsErrEnvNotFound(t *testing.T) {
	repo, _ := initTestDB(t)

	_, err := repo.GetEnvironment(context.Background(), 1)
	if err == nil {
		t.Fatal("Expected an error but got none!")
	}

	if err != environments.ErrNotFound {
		t.Fatalf("Expected an ErrEnvNotFound got: %s", err)
	}
}

func TestGetEnvironmentByName(t *testing.T) {
	repo, svcRepo := initTestDB(t)

	svc := svcFixture(t, svcRepo, "svc1")

	env1 := envFixture(t, repo, "env1", nil, svc.ID)
	env1.Service = svc.Name
	envFixture(t, repo, "env2", nil, svc.ID)

	env, err := repo.GetEnvironmentByName(context.Background(), svc.Name, env1.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(env1, env) {
		t.Fatalf("Got wrong environment expected %v got %v", env1, env)
	}
}

func TestGetEnvironmentByNameReturnsErrEnvNotFound(t *testing.T) {
	repo, _ := initTestDB(t)

	_, err := repo.GetEnvironmentByName(context.Background(), "service", "dev")
	if err == nil {
		t.Fatal("Expected an error but got none!")
	}

	if err != environments.ErrNotFound {
		t.Fatalf("Expected an ErrEnvNotFound got: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	repo, svcRepo := initTestDB(t)

	svc := svcFixture(t, svcRepo, "svc1")

	env1 := envFixture(t, repo, "env1", nil, svc.ID)
	env2 := envFixture(t, repo, "env2", &env1.ID, svc.ID)
	updated := environments.Environment{
		ID:           env2.ID,
		Name:         "updated",
		PromotesToID: nil,
		Sensitive:    true,
		ServiceID:    svc.ID,
		CreatedAt:    env2.CreatedAt,
	}

	result, err := repo.UpdateEnvironment(context.Background(), updated)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(updated, result) {
		t.Errorf("expected: %v got: %v", updated, result)
	}
}

func TestDelete(t *testing.T) {
	repo, svcRepo := initTestDB(t)

	svc := svcFixture(t, svcRepo, "svc1")

	env1 := envFixture(t, repo, "env1", nil, svc.ID)
	env2 := envFixture(t, repo, "env2", &env1.ID, svc.ID)

	err := repo.DeleteEnvironment(context.Background(), env2.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.GetEnvironment(context.Background(), env2.ID)

	if !errors.Is(err, environments.ErrNotFound) {
		t.Errorf("expected: %s got: %s", environments.ErrNotFound, err)
	}
}
