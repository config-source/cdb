package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/services"
)

func TestListConfigKeys(t *testing.T) {
	repo := &cdb.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "test",
			},
		},
		ConfigKeys: map[int]configkeys.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: configkeys.TypeString,
				ServiceID: 1,
				Service:   "test",
			},
			2: {
				ID:        2,
				Name:      "minReplicas",
				ValueType: configkeys.TypeInteger,
				ServiceID: 1,
				Service:   "test",
			},
			3: {
				ID:        3,
				Name:      "maxReplicas",
				ValueType: configkeys.TypeInteger,
				ServiceID: 1,
				Service:   "test",
			},
		},
	}

	_, mux, _ := testAPI(repo, true)
	req := httptest.NewRequest("GET", "/api/v1/config-keys", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var keys []configkeys.ConfigKey
	if err := json.NewDecoder(rr.Body).Decode(&keys); err != nil {
		t.Fatal(err)
	}

	if len(keys) != 3 {
		t.Fatalf("Expected 3 config keys got: %v", keys)
	}
}

func TestGetConfigKeyByID(t *testing.T) {
	repo := &cdb.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "test",
			},
		},
		ConfigKeys: map[int]configkeys.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: configkeys.TypeString,
				ServiceID: 1,
				Service:   "test",
			},
		},
	}

	_, mux, _ := testAPI(repo, true)
	req := httptest.NewRequest("GET", "/api/v1/config-keys/1/by-id/1", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetConfigKeyByIDNotFound(t *testing.T) {
	repo := &cdb.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "test",
			},
		},
		ConfigKeys: map[int]configkeys.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: configkeys.TypeString,
				ServiceID: 1,
				Service:   "test",
			},
		},
	}

	_, mux, _ := testAPI(repo, true)
	req := httptest.NewRequest("GET", "/api/v1/config-keys/1/by-id/2", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Fatalf("Expected status code 404 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestCreateConfigKey(t *testing.T) {
	repo := &cdb.TestRepository{}

	_, mux, _ := testAPI(repo, true)

	configKey := configkeys.ConfigKey{
		Name:      "owner",
		ValueType: configkeys.TypeString,
	}

	marshalled, err := json.Marshal(configKey)
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

	var created configkeys.ConfigKey
	err = json.NewDecoder(rr.Body).Decode(&created)
	if err != nil {
		t.Fatal(err)
	}

	if created.Name != configKey.Name {
		t.Fatalf("Expected name to be %s got: %s", configKey.Name, created.Name)
	}

	if created.CreatedAt.IsZero() {
		t.Fatal("Expected CreatedAt to not be zero value.")
	}

	if created.ID == 0 {
		t.Fatal("Expected ID to not be zero.")
	}
}

func TestGetConfigKeyByName(t *testing.T) {
	repo := &cdb.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "test",
			},
		},
		ConfigKeys: map[int]configkeys.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: configkeys.TypeString,
				ServiceID: 1,
				Service:   "test",
			},
		},
	}

	_, mux, _ := testAPI(repo, true)
	req := httptest.NewRequest("GET", "/api/v1/config-keys/1/by-name/owner", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetConfigKeyByNameNotFound(t *testing.T) {
	repo := &cdb.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "test",
			},
		},
		ConfigKeys: map[int]configkeys.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: configkeys.TypeString,
				ServiceID: 1,
				Service:   "test",
			},
		},
	}

	_, mux, _ := testAPI(repo, true)
	req := httptest.NewRequest("GET", "/api/v1/config-keys/1/by-name/minReplicas", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Fatalf("Expected status code 404 got: %d %s", rr.Code, rr.Body.String())
	}
}
