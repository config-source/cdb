package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/internal/repository"
)

func TestListConfigKeys(t *testing.T) {
	repo := &repository.TestRepository{
		ConfigKeys: map[int]cdb.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: cdb.TypeString,
			},
			2: {
				ID:        2,
				Name:      "minReplicas",
				ValueType: cdb.TypeInteger,
			},
			3: {
				ID:        3,
				Name:      "maxReplicas",
				ValueType: cdb.TypeInteger,
			},
		},
	}

	_, mux := testAPI(repo)
	req := httptest.NewRequest("GET", "/api/v1/config-keys", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var keys []cdb.ConfigKey
	if err := json.NewDecoder(rr.Body).Decode(&keys); err != nil {
		t.Fatal(err)
	}

	if len(keys) != 3 {
		t.Fatalf("Expected 3 config keys got: %v", keys)
	}
}

func TestGetConfigKeyByID(t *testing.T) {
	repo := &repository.TestRepository{
		ConfigKeys: map[int]cdb.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: cdb.TypeString,
			},
		},
	}

	_, mux := testAPI(repo)
	req := httptest.NewRequest("GET", "/api/v1/config-keys/by-id/1", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetConfigKeyByIDNotFound(t *testing.T) {
	repo := &repository.TestRepository{
		ConfigKeys: map[int]cdb.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: cdb.TypeString,
			},
		},
	}

	_, mux := testAPI(repo)
	req := httptest.NewRequest("GET", "/api/v1/config-keys/2", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Fatalf("Expected status code 404 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestCreateConfigKey(t *testing.T) {
	repo := &repository.TestRepository{}

	_, mux := testAPI(repo)

	env := cdb.ConfigKey{
		Name:      "owner",
		ValueType: cdb.TypeString,
	}

	marshalled, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/v1/config-keys", bytes.NewBuffer(marshalled))
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 201 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var created cdb.ConfigKey
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
