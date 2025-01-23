package services_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/config-source/cdb/pkg/postgresutils"
	"github.com/config-source/cdb/pkg/services"
	"github.com/rs/zerolog"
)

func initTestDB(t *testing.T) *services.Repository {
	t.Helper()

	pool := postgresutils.InitTestDB(t)
	repo := services.NewRepository(
		zerolog.New(nil).Level(zerolog.Disabled),
		pool,
	)

	return repo
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

func TestCreateService(t *testing.T) {
	repo := initTestDB(t)

	svc, err := repo.CreateService(context.Background(), services.Service{
		Name: "mat",
	})
	if err != nil {
		t.Fatal(err)
	}

	if svc.ID == 0 {
		t.Fatalf("Expected ID to be set got: %d", svc.ID)
	}

	if svc.Name != "mat" {
		t.Fatalf("Expected Name to be mat got: %s", svc.Name)
	}
}

func TestGetService(t *testing.T) {
	repo := initTestDB(t)

	svcFixture(t, repo, "svc1")
	svc2 := svcFixture(t, repo, "svc2")

	svc, err := repo.GetService(context.Background(), svc2.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(svc2, svc) {
		t.Fatalf("Got wrong service expected %v got %v", svc2, svc)
	}
}

func TestGetServiceReturnsErrsvcNotFound(t *testing.T) {
	repo := initTestDB(t)

	_, err := repo.GetService(context.Background(), 1)
	if err == nil {
		t.Fatal("Expected an error but got none!")
	}

	if err != services.ErrNotFound {
		t.Fatalf("Expected an ErrsvcNotFound got: %s", err)
	}
}

func TestGetServiceByName(t *testing.T) {
	repo := initTestDB(t)

	svc1 := svcFixture(t, repo, "svc1")
	svcFixture(t, repo, "svc2")

	svc, err := repo.GetServiceByName(context.Background(), svc1.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(svc1, svc) {
		t.Fatalf("Got wrong service expected %v got %v", svc1, svc)
	}
}

func TestGetServiceByNameReturnsErrsvcNotFound(t *testing.T) {
	repo := initTestDB(t)

	_, err := repo.GetServiceByName(context.Background(), "dev")
	if err == nil {
		t.Fatal("Expected an error but got none!")
	}

	if err != services.ErrNotFound {
		t.Fatalf("Expected an ErrsvcNotFound got: %s", err)
	}
}
