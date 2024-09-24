package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb/pkg/services"
	"github.com/config-source/cdb/pkg/testutils"
)

func TestGetServiceByName(t *testing.T) {
	repo := &testutils.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "production",
			},
		},
	}

	_, mux, _ := testAPI(repo, true)
	req := httptest.NewRequest("GET", "/api/v1/services/by-name/production", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetServiceByNameNotFound(t *testing.T) {
	repo := &testutils.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "production",
			},
		},
	}

	_, mux, _ := testAPI(repo, true)
	req := httptest.NewRequest("GET", "/api/v1/services/by-name/dev", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Fatalf("Expected status code 404 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetServiceByID(t *testing.T) {
	repo := &testutils.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "production",
			},
		},
	}

	_, mux, _ := testAPI(repo, true)
	req := httptest.NewRequest("GET", "/api/v1/services/by-id/1", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetServiceByIDNotFound(t *testing.T) {
	repo := &testutils.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "production",
			},
		},
	}

	_, mux, _ := testAPI(repo, true)
	req := httptest.NewRequest("GET", "/api/v1/services/by-id/2", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Fatalf("Expected status code 404 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestCreateService(t *testing.T) {
	repo := &testutils.TestRepository{}

	_, mux, _ := testAPI(repo, true)

	svc := services.Service{
		Name: "production",
	}

	marshalled, err := json.Marshal(svc)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/v1/services", bytes.NewBuffer(marshalled))
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 201 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var created services.Service
	err = json.NewDecoder(rr.Body).Decode(&created)
	if err != nil {
		t.Fatal(err)
	}

	if created.Name != svc.Name {
		t.Fatalf("Expected name to be %s got: %s", svc.Name, created.Name)
	}

	if created.CreatedAt.IsZero() {
		t.Fatal("Expected CreatedAt to not be zero value.")
	}

	if created.ID == 0 {
		t.Fatal("Expected ID to not be zero.")
	}
}
