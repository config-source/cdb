package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/services"
)

func createKeys(t *testing.T, tc TestContext, keys []configkeys.ConfigKey) []configkeys.ConfigKey {
	svc, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "keyHolder"})
	if err != nil {
		t.Fatal(err)
	}

	newKeys := make([]configkeys.ConfigKey, len(keys))
	for idx, key := range keys {
		if key.ServiceID == 0 {
			key.ServiceID = svc.ID
		}

		newKeys[idx], err = tc.keyRepo.CreateConfigKey(context.Background(), key)
		if err != nil {
			t.Fatal(err)
		}
	}

	return newKeys
}

func TestListConfigKeys(t *testing.T) {
	tc, mux := testAPI(t, true)

	createKeys(t, tc, []configkeys.ConfigKey{
		{
			Name:      "owner",
			ValueType: configkeys.TypeString,
		},
		{
			Name:      "minReplicas",
			ValueType: configkeys.TypeInteger,
		},
		{
			Name:      "maxReplicas",
			ValueType: configkeys.TypeInteger,
		},
	})

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
	tc, mux := testAPI(t, true)

	keys := createKeys(t, tc, []configkeys.ConfigKey{
		{
			Name:      "owner",
			ValueType: configkeys.TypeString,
		},
	})
	key := keys[0]

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/config-keys/by-id/%d", key.ID), nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var got configkeys.ConfigKey
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if got.ID != key.ID {
		t.Errorf("Expected ID %d got %d", key.ID, got.ID)
	}
}

func TestGetConfigKeyByIDNotFound(t *testing.T) {
	tc, mux := testAPI(t, true)

	svc, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}

	keys := createKeys(t, tc, []configkeys.ConfigKey{
		{
			Name:      "owner",
			ValueType: configkeys.TypeString,
			ServiceID: svc.ID,
		},
	})
	key := keys[0]

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/config-keys/%d/by-id/%d", svc.ID, key.ID+1), nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Fatalf("Expected status code 404 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestCreateConfigKey(t *testing.T) {
	tc, mux := testAPI(t, true)

	svc, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}

	configKey := configkeys.ConfigKey{
		Name:      "owner",
		ValueType: configkeys.TypeString,
		ServiceID: svc.ID,
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
		t.Errorf("Expected name to be %s got: %s", configKey.Name, created.Name)
	}

	if created.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to not be zero value.")
	}

	if created.ID == 0 {
		t.Error("Expected ID to not be zero.")
	}
}

func TestGetConfigKeyByName(t *testing.T) {
	tc, mux := testAPI(t, true)

	createKeys(t, tc, []configkeys.ConfigKey{
		{
			Name:      "owner",
			ValueType: configkeys.TypeString,
		},
	})

	req := httptest.NewRequest("GET", "/api/v1/config-keys/keyHolder/by-name/owner", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestGetConfigKeyByNameNotFound(t *testing.T) {
	tc, mux := testAPI(t, true)

	createKeys(t, tc, []configkeys.ConfigKey{
		{
			Name:      "owner",
			ValueType: configkeys.TypeString,
		},
	})

	req := httptest.NewRequest("GET", "/api/v1/config-keys/1/by-name/minReplicas", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 404 {
		t.Fatalf("Expected status code 404 got: %d %s", rr.Code, rr.Body.String())
	}
}

func TestListConfigKeysFilterByService(t *testing.T) {
	tc, mux := testAPI(t, true)

	svc1, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "test1"})
	if err != nil {
		t.Fatal(err)
	}

	svc2, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "test2"})
	if err != nil {
		t.Fatal(err)
	}

	svc3, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "test3"})
	if err != nil {
		t.Fatal(err)
	}

	createKeys(t, tc, []configkeys.ConfigKey{
		{
			Name:      "owner",
			ValueType: configkeys.TypeString,
			ServiceID: svc1.ID,
		},
		{
			Name:      "minReplicas",
			ValueType: configkeys.TypeInteger,
			ServiceID: svc2.ID,
		},
		{
			ID:        3,
			Name:      "maxReplicas",
			ValueType: configkeys.TypeInteger,
			ServiceID: svc3.ID,
		},
	})

	req := httptest.NewRequest(
		"GET",
		fmt.Sprintf("/api/v1/config-keys?service=%d&service=%d", svc2.ID, svc3.ID),
		nil,
	)
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
		if !slices.Contains([]int{svc2.ID, svc3.ID}, key.ServiceID) {
			t.Fatalf("Expected ServiceID to be %d or %d got: %v", svc2.ID, svc3.ID, keys[0])
		}
	}
}
