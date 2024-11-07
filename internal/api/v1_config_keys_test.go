package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/services"
	"github.com/config-source/cdb/pkg/testutils"
)

func TestListConfigKeys(t *testing.T) {
	repo := &testutils.TestRepository{
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
	repo := &testutils.TestRepository{
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
	req := httptest.NewRequest("GET", "/api/v1/config-keys/by-id/1", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetConfigKeyByIDNotFound(t *testing.T) {
	repo := &testutils.TestRepository{
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
	repo := &testutils.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "test",
			},
		},
	}

	_, mux, _ := testAPI(repo, true)

	configKey := configkeys.ConfigKey{
		Name:      "owner",
		ValueType: configkeys.TypeString,
		ServiceID: 1,
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
	repo := &testutils.TestRepository{
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
	req := httptest.NewRequest("GET", "/api/v1/config-keys/test/by-name/owner", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetConfigKeyByNameNotFound(t *testing.T) {
	repo := &testutils.TestRepository{
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

func TestListConfigKeysFilterByService(t *testing.T) {
	repo := &testutils.TestRepository{
		Services: map[int]services.Service{
			1: {
				ID:   1,
				Name: "test",
			},
			2: {
				ID:   2,
				Name: "test2",
			},
			3: {
				ID:   3,
				Name: "test3",
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
				ServiceID: 2,
				Service:   "test2",
			},
			3: {
				ID:        3,
				Name:      "maxReplicas",
				ValueType: configkeys.TypeInteger,
				ServiceID: 3,
				Service:   "test3",
			},
		},
	}

	_, mux, _ := testAPI(repo, true)
	req := httptest.NewRequest("GET", "/api/v1/config-keys?service=2&service=3", nil)
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

	if len(keys) != 2 {
		t.Fatalf("Expected 2 config keys got: %v", keys)
	}

	for _, key := range keys {
		if !slices.Contains([]int{2, 3}, key.ServiceID) {
			t.Fatalf("Expected ServiceID to be 2 or 3 got: %v", keys[0])
		}
	}
}
