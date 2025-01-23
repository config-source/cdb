package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb/pkg/services"
)

func TestGetServiceByName(t *testing.T) {
	tc, mux := testAPI(t, true)
	svc, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "production"})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/api/v1/services/by-name/production", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var got services.Service
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if got.ID != svc.ID {
		t.Errorf("Expected ID %d got %d", svc.ID, got.ID)
	}
}

func TestGetServiceByNameNotFound(t *testing.T) {
	tc, mux := testAPI(t, true)
	_, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "production"})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/api/v1/services/by-name/dev", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Fatalf("Expected status code 404 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetServiceByID(t *testing.T) {
	tc, mux := testAPI(t, true)
	svc, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "production"})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/services/by-id/%d", svc.ID), nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetServiceByIDNotFound(t *testing.T) {
	tc, mux := testAPI(t, true)
	_, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "production"})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/api/v1/services/by-id/2", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Fatalf("Expected status code 404 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestCreateService(t *testing.T) {
	tc, mux := testAPI(t, true)

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

	_, err = tc.serviceRepo.GetService(context.Background(), created.ID)
	if err != nil {
		t.Fatal(err)
	}
}
