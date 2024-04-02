package api

import (
	"bytes"
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
