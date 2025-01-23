package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/services"
)

func TestGetEnvironmentByName(t *testing.T) {
	tc, mux := testAPI(t, true)

	svc, err := tc.serviceRepo.CreateService(
		context.Background(),
		services.Service{Name: "test"},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tc.environmentRepo.CreateEnvironment(
		context.Background(),
		environments.Environment{Name: "production", ServiceID: svc.ID},
	)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/api/v1/environments/test/by-name/production", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetEnvironmentByNameNotFound(t *testing.T) {
	tc, mux := testAPI(t, true)

	svc, err := tc.serviceRepo.CreateService(
		context.Background(),
		services.Service{Name: "test"},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tc.environmentRepo.CreateEnvironment(
		context.Background(),
		environments.Environment{Name: "production", ServiceID: svc.ID},
	)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/api/v1/environments/test/by-name/dev", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Fatalf("Expected status code 404 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetEnvironmentByID(t *testing.T) {
	tc, mux := testAPI(t, true)

	svc, err := tc.serviceRepo.CreateService(
		context.Background(),
		services.Service{Name: "test"},
	)
	if err != nil {
		t.Fatal(err)
	}

	env, err := tc.environmentRepo.CreateEnvironment(
		context.Background(),
		environments.Environment{Name: "production", ServiceID: svc.ID},
	)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/environments/by-id/%d", env.ID), nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetEnvironmentByIDNotFound(t *testing.T) {
	tc, mux := testAPI(t, true)

	svc, err := tc.serviceRepo.CreateService(
		context.Background(),
		services.Service{Name: "test"},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tc.environmentRepo.CreateEnvironment(
		context.Background(),
		environments.Environment{Name: "production", ServiceID: svc.ID},
	)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/api/v1/environments/by-id/2", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Fatalf("Expected status code 404 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestCreateEnvironment(t *testing.T) {
	tc, mux := testAPI(t, true)

	svc, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}

	env := environments.Environment{
		ServiceID: svc.ID,
		Name:      "production",
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

	var created environments.Environment
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
