package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/internal/repository"
)

func TestGetEnvironmentByName(t *testing.T) {
	repo := &repository.TestRepository{
		Environments: map[int]cdb.Environment{
			1: {
				ID:   1,
				Name: "production",
			},
		},
	}

	_, mux := testAPI(repo)
	req := httptest.NewRequest("GET", "/api/v1/environments/by-name/production", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetEnvironmentByNameNotFound(t *testing.T) {
	repo := &repository.TestRepository{
		Environments: map[int]cdb.Environment{
			1: {
				ID:   1,
				Name: "production",
			},
		},
	}

	_, mux := testAPI(repo)
	req := httptest.NewRequest("GET", "/api/v1/environments/by-name/dev", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Fatalf("Expected status code 404 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetEnvironmentByID(t *testing.T) {
	repo := &repository.TestRepository{
		Environments: map[int]cdb.Environment{
			1: {
				ID:   1,
				Name: "production",
			},
		},
	}

	_, mux := testAPI(repo)
	req := httptest.NewRequest("GET", "/api/v1/environments/by-id/1", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetEnvironmentByIDNotFound(t *testing.T) {
	repo := &repository.TestRepository{
		Environments: map[int]cdb.Environment{
			1: {
				ID:   1,
				Name: "production",
			},
		},
	}

	_, mux := testAPI(repo)
	req := httptest.NewRequest("GET", "/api/v1/environments/by-id/2", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Fatalf("Expected status code 404 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestCreateEnvironment(t *testing.T) {
	repo := &repository.TestRepository{}

	_, mux := testAPI(repo)

	env := cdb.Environment{
		Name: "production",
	}

	marshalled, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/v1/environments", bytes.NewBuffer(marshalled))
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 201 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var created cdb.Environment
	err = json.NewDecoder(rr.Body).Decode(&created)
	if err != nil {
		t.Fatal(err)
	}

	if created.Name != env.Name {
		t.Fatalf("Expected name to be %s got: %s", env.Name, created.Name)
	}

	if created.CreatedAt.IsZero() {
		t.Fatal("Expected CreatedAt to not be zero value.")
	}

	if created.ID == 0 {
		t.Fatal("Expected ID to not be zero.")
	}
}
