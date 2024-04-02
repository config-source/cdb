package api

import (
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb/internal/repository"
)

func TestHealthCheckSuccess(t *testing.T) {
	repo := &repository.TestRepository{
		IsHealthy: true,
	}

	_, mux := testAPI(repo)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected 200 status code got: %d", rr.Code)
	}
}

func TestHealthCheckFailure(t *testing.T) {
	repo := &repository.TestRepository{
		IsHealthy: false,
	}

	_, mux := testAPI(repo)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code == 200 {
		t.Fatalf("Expected non-200 status code got: %d", rr.Code)
	}
}
