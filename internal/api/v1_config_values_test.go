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
	"github.com/config-source/cdb/pkg/services"
)

func TestGetConfiguration(t *testing.T) {
	tc, mux := testAPI(t, true)

	svc, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}

	production, err := tc.environmentRepo.CreateEnvironment(context.Background(), environments.Environment{Name: "production", ServiceID: svc.ID})
	if err != nil {
		t.Fatal(err)
	}

	staging, err := tc.environmentRepo.CreateEnvironment(
		context.Background(),
		environments.Environment{
			Name:         "staging",
			ServiceID:    svc.ID,
			PromotesToID: &production.ID,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	keys := createKeys(t, tc, []configkeys.ConfigKey{
		{
			Name:      "owner",
			ValueType: configkeys.TypeString,
			ServiceID: svc.ID,
		},
		{
			Name:      "maxReplicas",
			ValueType: configkeys.TypeInteger,
			ServiceID: svc.ID,
		},
	})

	fixtures := []*configvalues.ConfigValue{
		configvalues.NewString(production.ID, keys[0].ID, "SRE"),
		configvalues.NewInt(production.ID, keys[1].ID, 100),
		configvalues.NewInt(staging.ID, keys[1].ID, 10),
	}

	for _, cv := range fixtures {
		_, err := tc.valueRepo.CreateConfigValue(context.Background(), cv)
		if err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest("GET", "/api/v1/config-values/2", nil)
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
	tc, mux := testAPI(t, true)

	svc, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}

	production, err := tc.environmentRepo.CreateEnvironment(context.Background(), environments.Environment{Name: "production", ServiceID: svc.ID})
	if err != nil {
		t.Fatal(err)
	}

	key := createKeys(t, tc, []configkeys.ConfigKey{
		{
			Name:      "owner",
			ValueType: configkeys.TypeString,
			ServiceID: svc.ID,
		},
	})[0]

	val := "test"
	value := configvalues.ConfigValue{
		Name:          "owner",
		ValueType:     configkeys.TypeString,
		EnvironmentID: production.ID,
		ConfigKeyID:   key.ID,
		StrValue:      &val,
	}

	marshalled, err := json.Marshal(value)
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

	if created.Name != value.Name {
		t.Errorf("Expected name to be '%s' got: '%s'", value.Name, created.Name)
	}

	if created.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to not be zero value.")
	}

	if created.ID == 0 {
		t.Error("Expected ID to not be zero.")
	}
}

func TestGetConfigurationByKey(t *testing.T) {
	tc, mux := testAPI(t, true)

	svc, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}

	production, err := tc.environmentRepo.CreateEnvironment(context.Background(), environments.Environment{Name: "production", ServiceID: svc.ID})
	if err != nil {
		t.Fatal(err)
	}

	staging, err := tc.environmentRepo.CreateEnvironment(
		context.Background(),
		environments.Environment{
			Name:         "staging",
			ServiceID:    svc.ID,
			PromotesToID: &production.ID,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	keys := createKeys(t, tc, []configkeys.ConfigKey{
		{
			Name:      "owner",
			ValueType: configkeys.TypeString,
			ServiceID: svc.ID,
		},
		{
			Name:      "maxReplicas",
			ValueType: configkeys.TypeInteger,
			ServiceID: svc.ID,
		},
	})

	fixtures := []*configvalues.ConfigValue{
		configvalues.NewString(production.ID, keys[0].ID, "SRE"),
		configvalues.NewInt(production.ID, keys[1].ID, 100),
		configvalues.NewInt(staging.ID, keys[1].ID, 10),
	}

	for _, cv := range fixtures {
		_, err := tc.valueRepo.CreateConfigValue(context.Background(), cv)
		if err != nil {
			t.Fatal(err)
		}
	}

	req := httptest.NewRequest("GET", "/api/v1/config-values/2/owner", nil)
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
	tc, mux := testAPI(t, true)

	svc, err := tc.serviceRepo.CreateService(context.Background(), services.Service{Name: "test"})
	if err != nil {
		t.Fatal(err)
	}

	production, err := tc.environmentRepo.CreateEnvironment(context.Background(), environments.Environment{Name: "production", ServiceID: svc.ID})
	if err != nil {
		t.Fatal(err)
	}

	staging, err := tc.environmentRepo.CreateEnvironment(
		context.Background(),
		environments.Environment{
			Name:         "staging",
			ServiceID:    svc.ID,
			PromotesToID: &production.ID,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	keys := createKeys(t, tc, []configkeys.ConfigKey{
		{
			Name:      "owner",
			ValueType: configkeys.TypeString,
			ServiceID: svc.ID,
		},
		{
			Name:      "maxReplicas",
			ValueType: configkeys.TypeInteger,
			ServiceID: svc.ID,
		},
	})

	fixtures := []*configvalues.ConfigValue{
		configvalues.NewString(production.ID, keys[0].ID, "SRE"),
		configvalues.NewInt(production.ID, keys[1].ID, 100),
		configvalues.NewInt(staging.ID, keys[1].ID, 10),
	}

	for _, cv := range fixtures {
		_, err := tc.valueRepo.CreateConfigValue(context.Background(), cv)
		if err != nil {
			t.Fatal(err)
		}
	}

	val := 10
	env := configvalues.ConfigValue{
		ValueType: configkeys.TypeInteger,
		IntValue:  &val,
	}

	marshalled, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/v1/config-values/2/minReplicas", bytes.NewBuffer(marshalled))
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
