package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/configvalues"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/testutils"
)

func TestGetConfiguration(t *testing.T) {
	promotesToID := 1
	repo := &testutils.TestRepository{
		Environments: map[int]environments.Environment{
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
		ConfigKeys: map[int]configkeys.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: configkeys.TypeString,
			},
			2: {
				ID:        2,
				Name:      "maxReplicas",
				ValueType: configkeys.TypeInteger,
			},
		},
	}

	fixtures := []*configvalues.ConfigValue{
		configvalues.NewString(1, 1, "SRE"),
		configvalues.NewInt(1, 2, 100),
		configvalues.NewInt(2, 2, 10),
	}

	for _, cv := range fixtures {
		_, err := repo.CreateConfigValue(context.Background(), cv)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, mux, _ := testAPI(repo, true)
	req := httptest.NewRequest("GET", "/api/v1/config-values/staging", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var values []configvalues.ConfigValue
	if err := json.NewDecoder(rr.Body).Decode(&values); err != nil {
		t.Fatal(err)
	}

	if len(values) != 2 {
		t.Fatalf("Expected 2 config values got: %v", values)
	}

	owner := values[1]
	if owner.ValueType != configkeys.TypeString {
		t.Fatalf("Expected configkeys.TypeString (%d) config got: %v", configkeys.TypeString, owner)
	}

	if owner.Value().(string) != "SRE" {
		t.Fatalf("Expected \"SRE\" got: %v", owner)
	}

	maxReplicas := values[0]
	if maxReplicas.ValueType != configkeys.TypeInteger {
		t.Fatalf("Expected configkeys.TypeInteger (%d) config got: %v", configkeys.TypeInteger, maxReplicas)
	}

	if maxReplicas.Value().(int) != 10 {
		t.Fatalf("Expected 10 got: %v", maxReplicas)
	}
}

func TestCreateConfigValue(t *testing.T) {
	repo := &testutils.TestRepository{
		Environments: map[int]environments.Environment{
			1: {
				ID:   1,
				Name: "production",
			},
		},
		ConfigKeys: map[int]configkeys.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: configkeys.TypeString,
			},
		},
	}

	_, mux, _ := testAPI(repo, true)

	val := "test"
	env := configvalues.ConfigValue{
		Name:          "owner",
		ValueType:     configkeys.TypeString,
		EnvironmentID: 1,
		ConfigKeyID:   1,
		StrValue:      &val,
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

	var created configvalues.ConfigValue
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
	repo := &testutils.TestRepository{
		Environments: map[int]environments.Environment{
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
		ConfigKeys: map[int]configkeys.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: configkeys.TypeString,
			},
			2: {
				ID:        2,
				Name:      "maxReplicas",
				ValueType: configkeys.TypeInteger,
			},
		},
	}

	fixtures := []*configvalues.ConfigValue{
		configvalues.NewString(1, 1, "SRE"),
		configvalues.NewInt(1, 2, 100),
		configvalues.NewInt(2, 2, 10),
	}

	for _, cv := range fixtures {
		_, err := repo.CreateConfigValue(context.Background(), cv)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, mux, _ := testAPI(repo, true)
	req := httptest.NewRequest("GET", "/api/v1/config-values/staging/owner", nil)
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var owner configvalues.ConfigValue
	if err := json.NewDecoder(rr.Body).Decode(&owner); err != nil {
		t.Fatal(err)
	}

	if owner.ValueType != configkeys.TypeString {
		t.Fatalf("Expected configkeys.TypeString (%d) config got: %v", configkeys.TypeString, owner)
	}

	if owner.Value().(string) != "SRE" {
		t.Fatalf("Expected \"SRE\" got: %v", owner)
	}
}

func TestSetConfigurationByKey(t *testing.T) {
	promotesToID := 1
	repo := &testutils.TestRepository{
		Environments: map[int]environments.Environment{
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
		ConfigKeys: map[int]configkeys.ConfigKey{
			1: {
				ID:        1,
				Name:      "owner",
				ValueType: configkeys.TypeString,
			},
			2: {
				ID:        2,
				Name:      "maxReplicas",
				ValueType: configkeys.TypeInteger,
			},
		},
	}

	fixtures := []*configvalues.ConfigValue{
		configvalues.NewString(1, 1, "SRE"),
		configvalues.NewInt(1, 2, 100),
		configvalues.NewInt(2, 2, 10),
	}

	for _, cv := range fixtures {
		_, err := repo.CreateConfigValue(context.Background(), cv)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, mux, _ := testAPI(repo, true)

	val := 10
	env := configvalues.ConfigValue{
		ValueType: configkeys.TypeInteger,
		IntValue:  &val,
	}

	marshalled, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/v1/config-values/staging/minReplicas", bytes.NewBuffer(marshalled))
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var maxReplicas configvalues.ConfigValue
	if err := json.NewDecoder(rr.Body).Decode(&maxReplicas); err != nil {
		t.Fatal(err)
	}

	if maxReplicas.ValueType != configkeys.TypeInteger {
		t.Fatalf("Expected configkeys.TypeInteger (%d) config got: %v", configkeys.TypeInteger, maxReplicas)
	}

	if maxReplicas.Value().(int) != 10 {
		t.Fatalf("Expected 10 got: %v", maxReplicas)
	}
}
