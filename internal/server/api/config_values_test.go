package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/internal/repository"
)

func TestGetConfiguration(t *testing.T) {
	promotesToID := 1
	repo := &repository.TestRepository{
		Environments: map[int]cdb.Environment{
			1: {
				ID:   1,
				Name: "production",
			},
			2: {
				ID:           2,
				Name:         "staging",
				PromotesToID: &promotesToID,
			},
		},
		ConfigKeys: map[int]cdb.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: cdb.TypeString,
			},
			2: {
				ID:        2,
				Name:      "maxReplicas",
				ValueType: cdb.TypeInteger,
			},
		},
	}

	fixtures := []cdb.ConfigValue{
		cdb.NewStringConfigValue(1, 1, "SRE"),
		cdb.NewIntConfigValue(1, 2, 100),
		cdb.NewIntConfigValue(2, 2, 10),
	}

	for _, cv := range fixtures {
		_, err := repo.CreateConfigValue(context.Background(), cv)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, mux := testAPI(repo)
	req := httptest.NewRequest("GET", "/api/v1/config-values/staging", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var values []cdb.ConfigValue
	if err := json.NewDecoder(rr.Body).Decode(&values); err != nil {
		t.Fatal(err)
	}

	if len(values) != 2 {
		t.Fatalf("Expected 2 config values got: %v", values)
	}

	owner := values[1]
	if owner.ValueType != cdb.TypeString {
		t.Fatalf("Expected cdb.TypeString (%d) config got: %v", cdb.TypeString, owner)
	}

	if owner.Value().(string) != "SRE" {
		t.Fatalf("Expected \"SRE\" got: %v", owner)
	}

	maxReplicas := values[0]
	if maxReplicas.ValueType != cdb.TypeInteger {
		t.Fatalf("Expected cdb.TypeInteger (%d) config got: %v", cdb.TypeInteger, maxReplicas)
	}

	if maxReplicas.Value().(int) != 10 {
		t.Fatalf("Expected 10 got: %v", maxReplicas)
	}
}

func TestCreateConfigValue(t *testing.T) {
	repo := &repository.TestRepository{
		Environments: map[int]cdb.Environment{
			1: {
				ID:   1,
				Name: "production",
			},
		},
		ConfigKeys: map[int]cdb.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: cdb.TypeString,
			},
		},
	}

	_, mux := testAPI(repo)

	env := cdb.ConfigValue{
		Name:          "owner",
		ValueType:     cdb.TypeString,
		EnvironmentID: 1,
		ConfigKeyID:   1,
	}

	marshalled, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/v1/config-values", bytes.NewBuffer(marshalled))
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 201 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var created cdb.ConfigValue
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

func TestGetConfigurationByKey(t *testing.T) {
	promotesToID := 1
	repo := &repository.TestRepository{
		Environments: map[int]cdb.Environment{
			1: {
				ID:   1,
				Name: "production",
			},
			2: {
				ID:           2,
				Name:         "staging",
				PromotesToID: &promotesToID,
			},
		},
		ConfigKeys: map[int]cdb.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: cdb.TypeString,
			},
			2: {
				ID:        2,
				Name:      "maxReplicas",
				ValueType: cdb.TypeInteger,
			},
		},
	}

	fixtures := []cdb.ConfigValue{
		cdb.NewStringConfigValue(1, 1, "SRE"),
		cdb.NewIntConfigValue(1, 2, 100),
		cdb.NewIntConfigValue(2, 2, 10),
	}

	for _, cv := range fixtures {
		_, err := repo.CreateConfigValue(context.Background(), cv)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, mux := testAPI(repo)
	req := httptest.NewRequest("GET", "/api/v1/config-values/staging/owner", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var owner cdb.ConfigValue
	if err := json.NewDecoder(rr.Body).Decode(&owner); err != nil {
		t.Fatal(err)
	}

	if owner.ValueType != cdb.TypeString {
		t.Fatalf("Expected cdb.TypeString (%d) config got: %v", cdb.TypeString, owner)
	}

	if owner.Value().(string) != "SRE" {
		t.Fatalf("Expected \"SRE\" got: %v", owner)
	}
}
