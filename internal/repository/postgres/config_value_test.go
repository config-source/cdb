package postgres_test

import (
	"context"
	_ "embed"
	"errors"
	"reflect"
	"testing"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/internal/repository/postgres"
)

func TestCreateConfigValue(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	env := envFixture(t, repo, "cdb", nil)
	key := configKeyFixture(t, repo, "owner", cdb.TypeString, true)

	val := "test"
	cv, err := repo.CreateConfigValue(
		context.Background(),
		cdb.NewStringConfigValue(env.ID, key.ID, val),
	)
	if err != nil {
		t.Fatal(err)
	}

	if cv.ID == 0 {
		t.Fatalf("Expected ID to be set got: %d", cv.ID)
	}

	if *cv.StrValue != val {
		t.Fatalf("Expected string value %s got %s", val, *cv.StrValue)
	}

	if cv.IntValue != nil {
		t.Fatalf("Expected IntValue to be nil got: %+v", cv.IntValue)
	}

	if cv.FloatValue != nil {
		t.Fatalf("Expected FloatValue to be nil got: %+v", cv.FloatValue)
	}

	if cv.BoolValue != nil {
		t.Fatalf("Expected BoolValue to be nil got: %+v", cv.BoolValue)
	}

	retrieved := cv.Value()
	switch retrieved.(type) {
	case string:
		if retrieved != val {
			t.Fatalf("Expected %s got %s", val, retrieved)
		}
	default:
		t.Fatalf("Expected a string got %+v", retrieved)
	}
}

func TestUpdateConfigValue(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	env := envFixture(t, repo, "cdb", nil)
	key := configKeyFixture(t, repo, "owner", cdb.TypeString, true)
	cv := createConfigValue(t, repo, cdb.NewStringConfigValue(
		env.ID,
		key.ID,
		"test",
	))
	cv.SetStrValue("updated")

	var err error
	cv, err = repo.UpdateConfigurationValue(context.Background(), cv)
	if err != nil {
		t.Fatal(err)
	}

	if cv.ID == 0 {
		t.Fatalf("Expected ID to be set got: %d", cv.ID)
	}

	if *cv.StrValue != "updated" {
		t.Fatalf("Expected string value %s got %s", "updated", *cv.StrValue)
	}

	if cv.IntValue != nil {
		t.Fatalf("Expected IntValue to be nil got: %+v", cv.IntValue)
	}

	if cv.FloatValue != nil {
		t.Fatalf("Expected FloatValue to be nil got: %+v", cv.FloatValue)
	}

	if cv.BoolValue != nil {
		t.Fatalf("Expected BoolValue to be nil got: %+v", cv.BoolValue)
	}

	retrieved := cv.Value()
	switch retrieved.(type) {
	case string:
		if retrieved != "updated" {
			t.Fatalf("Expected %s got %s", "updated", retrieved)
		}
	default:
		t.Fatalf("Expected a string got %+v", retrieved)
	}
}

func TestUpdateConfigValueReturnsErrConfigValueNotFound(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	env := envFixture(t, repo, "cdb", nil)
	key := configKeyFixture(t, repo, "owner", cdb.TypeString, true)
	cv := createConfigValue(t, repo, cdb.NewStringConfigValue(
		env.ID,
		key.ID,
		"test",
	))
	cv.ID = cv.ID + 1

	var err error
	_, err = repo.UpdateConfigurationValue(context.Background(), cv)
	expectedError := cdb.ErrConfigValueNotFound
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected %s Got %s", expectedError, err)
	}
}

func TestUpdateConfigValueReturnsErrConfigKeyNotFound(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	env := envFixture(t, repo, "cdb", nil)
	key := configKeyFixture(t, repo, "owner", cdb.TypeString, true)
	cv := createConfigValue(t, repo, cdb.NewStringConfigValue(
		env.ID,
		key.ID,
		"test",
	))
	cv.ConfigKeyID = cv.ConfigKeyID + 1

	var err error
	_, err = repo.UpdateConfigurationValue(context.Background(), cv)
	expectedError := cdb.ErrConfigKeyNotFound
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected %s Got %s", expectedError, err)
	}
}

func TestUpdateConfigValueReturnsErrEnvironmentNotFound(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	env := envFixture(t, repo, "cdb", nil)
	key := configKeyFixture(t, repo, "owner", cdb.TypeString, true)
	cv := createConfigValue(t, repo, cdb.NewStringConfigValue(
		env.ID,
		key.ID,
		"test",
	))
	cv.EnvironmentID = cv.EnvironmentID + 1

	var err error
	_, err = repo.UpdateConfigurationValue(context.Background(), cv)
	expectedError := cdb.ErrEnvNotFound
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected %s Got %s", expectedError, err)
	}
}

func TestCreateIntConfigValue(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	env := envFixture(t, repo, "cdb", nil)
	key := configKeyFixture(t, repo, "owner", cdb.TypeInteger, true)

	val := "test"
	cv, err := repo.CreateConfigValue(
		context.Background(),
		cdb.NewStringConfigValue(env.ID, key.ID, val),
	)
	if err != nil {
		t.Fatal(err)
	}

	if cv.ID == 0 {
		t.Fatalf("Expected ID to be set got: %d", cv.ID)
	}

	if *cv.StrValue != val {
		t.Fatalf("Expected string value %s got %s", val, *cv.StrValue)
	}

	if cv.IntValue != nil {
		t.Fatalf("Expected IntValue to be nil got: %+v", cv.IntValue)
	}

	if cv.FloatValue != nil {
		t.Fatalf("Expected FloatValue to be nil got: %+v", cv.FloatValue)
	}

	if cv.BoolValue != nil {
		t.Fatalf("Expected BoolValue to be nil got: %+v", cv.BoolValue)
	}

	retrieved := cv.Value()
	switch retrieved.(type) {
	case string:
		if retrieved != val {
			t.Fatalf("Expected %s got %s", val, retrieved)
		}
	default:
		t.Fatalf("Expected a string got %+v", retrieved)
	}
}

func TestGetConfigValue(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	env := envFixture(t, repo, "cdb", nil)
	key := configKeyFixture(t, repo, "owner", cdb.TypeString, true)
	secondKey := configKeyFixture(t, repo, "secondKey", cdb.TypeString, true)

	_, err := repo.CreateConfigValue(context.Background(), cdb.NewStringConfigValue(
		env.ID,
		secondKey.ID,
		"test",
	))
	if err != nil {
		t.Fatal(err)
	}

	cv, err := repo.CreateConfigValue(context.Background(), cdb.NewStringConfigValue(
		env.ID,
		key.ID,
		"test",
	))
	if err != nil {
		t.Fatal(err)
	}

	retrieved, err := repo.GetConfigurationValue(context.Background(), env.Name, key.Name)
	if err != nil {
		t.Fatal(err)
	}

	cv.Name = key.Name
	cv.ValueType = key.ValueType

	if !reflect.DeepEqual(retrieved, cv) {
		t.Fatalf("Got wrong configValueironment expected %+v got %+v", cv, retrieved)
	}
}

func TestGetConfigValueInheritsValues(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	production := envFixture(t, repo, "production", nil)
	staging := envFixture(t, repo, "staging", &production.ID)
	dev := envFixture(t, repo, "dev", &staging.ID)

	owner := configKeyFixture(t, repo, "owner", cdb.TypeString, true)
	minReplicas := configKeyFixture(t, repo, "minReplicas", cdb.TypeInteger, true)
	maxReplicas := configKeyFixture(t, repo, "maxReplicas", cdb.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, cdb.NewIntConfigValue(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, maxReplicas.ID, 100))

	setDirectly := createConfigValue(t, repo, cdb.NewIntConfigValue(dev.ID, minReplicas.ID, 1))
	stagingInherited := createInheritedConfigValue(t, repo, staging.Name, cdb.NewIntConfigValue(staging.ID, maxReplicas.ID, 50))
	productionInherited := createInheritedConfigValue(t, repo, production.Name, cdb.NewStringConfigValue(production.ID, owner.ID, "SRE"))

	setDirectlyValue, err := repo.GetConfigurationValue(context.Background(), dev.Name, setDirectly.Name)
	if err != nil {
		t.Fatal(err)
	}

	if setDirectly.ID != setDirectlyValue.ID {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", setDirectly, setDirectlyValue)
	}

	stagingValue, err := repo.GetConfigurationValue(context.Background(), dev.Name, stagingInherited.Name)
	if err != nil {
		t.Fatal(err)
	}

	if stagingInherited.ID != stagingValue.ID {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", stagingInherited, stagingValue)
	}

	if !stagingValue.Inherited {
		t.Fatalf("Expected staging inherited value to be marked as such: %t\n", stagingValue.Inherited)
	}

	productionValue, err := repo.GetConfigurationValue(context.Background(), dev.Name, productionInherited.Name)
	if err != nil {
		t.Fatal(err)
	}

	if productionInherited.ID != productionValue.ID {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", productionInherited, productionValue)
	}

	if !productionValue.Inherited {
		t.Fatalf("Expected production inherited value to be marked as such: %t\n", productionValue.Inherited)
	}
}

func TestGetConfigValueReturnsCorrectErrorForValueNotFound(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	production := envFixture(t, repo, "production", nil)
	staging := envFixture(t, repo, "staging", &production.ID)
	dev := envFixture(t, repo, "dev", &staging.ID)

	owner := configKeyFixture(t, repo, "owner", cdb.TypeString, true)
	minReplicas := configKeyFixture(t, repo, "minReplicas", cdb.TypeInteger, true)
	maxReplicas := configKeyFixture(t, repo, "maxReplicas", cdb.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, cdb.NewIntConfigValue(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, maxReplicas.ID, 100))
	createConfigValue(t, repo, cdb.NewIntConfigValue(dev.ID, minReplicas.ID, 1))
	createInheritedConfigValue(t, repo, staging.Name, cdb.NewIntConfigValue(staging.ID, maxReplicas.ID, 50))
	createInheritedConfigValue(t, repo, production.Name, cdb.NewStringConfigValue(production.ID, owner.ID, "SRE"))

	_, err := repo.GetConfigurationValue(context.Background(), dev.Name, "notfound")
	if err != cdb.ErrConfigValueNotFound {
		t.Fatalf("Expected: %s Got: %s", cdb.ErrConfigValueNotFound, err)
	}
}

func TestGetConfigValueReturnsCorrectErrorForEnvNotFound(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	production := envFixture(t, repo, "production", nil)
	staging := envFixture(t, repo, "staging", &production.ID)
	dev := envFixture(t, repo, "dev", &staging.ID)

	owner := configKeyFixture(t, repo, "owner", cdb.TypeString, true)
	minReplicas := configKeyFixture(t, repo, "minReplicas", cdb.TypeInteger, true)
	maxReplicas := configKeyFixture(t, repo, "maxReplicas", cdb.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, cdb.NewIntConfigValue(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, maxReplicas.ID, 100))
	createConfigValue(t, repo, cdb.NewIntConfigValue(dev.ID, minReplicas.ID, 1))
	createInheritedConfigValue(t, repo, staging.Name, cdb.NewIntConfigValue(staging.ID, maxReplicas.ID, 50))
	createInheritedConfigValue(t, repo, production.Name, cdb.NewStringConfigValue(production.ID, owner.ID, "SRE"))

	_, err := repo.GetConfigurationValue(context.Background(), "notFound", "notfound")
	if !errors.Is(err, cdb.ErrConfigValueNotFound) {
		t.Fatalf("Expected: %s Got: %s", cdb.ErrEnvNotFound, err)
	}
}

func createConfigValue(t *testing.T, repo *postgres.Repository, cv *cdb.ConfigValue) *cdb.ConfigValue {
	created, err := repo.CreateConfigValue(context.Background(), cv)
	if err != nil {
		t.Fatal(err)
	}

	// Populate the config key fields by retrieving it as insert doesn't return
	// those.
	retrieved, err := repo.GetConfigurationValueByID(context.Background(), created.ID)
	if err != nil {
		t.Fatal(err)
	}

	return retrieved
}

func createInheritedConfigValue(t *testing.T, repo *postgres.Repository, parentName string, cv *cdb.ConfigValue) *cdb.ConfigValue {
	created := createConfigValue(t, repo, cv)
	created.Inherited = true
	created.InheritedFrom = parentName
	return created
}

func TestGetConfiguration(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	production := envFixture(t, repo, "production", nil)
	staging := envFixture(t, repo, "staging", &production.ID)
	dev := envFixture(t, repo, "dev", &staging.ID)

	owner := configKeyFixture(t, repo, "owner", cdb.TypeString, true)
	minReplicas := configKeyFixture(t, repo, "minReplicas", cdb.TypeInteger, true)
	maxReplicas := configKeyFixture(t, repo, "maxReplicas", cdb.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, cdb.NewIntConfigValue(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, maxReplicas.ID, 100))

	expectedValues := []cdb.ConfigValue{
		*createConfigValue(t, repo, cdb.NewIntConfigValue(dev.ID, minReplicas.ID, 1)),
		*createInheritedConfigValue(t, repo, staging.Name, cdb.NewIntConfigValue(staging.ID, maxReplicas.ID, 50)),
		*createInheritedConfigValue(t, repo, production.Name, cdb.NewStringConfigValue(production.ID, owner.ID, "SRE")),
	}

	retrieved, err := repo.GetConfiguration(context.Background(), dev.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedValues, retrieved) {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", expectedValues, retrieved)
	}
}

func TestGetConfigurationDoesntPropagateKeysWhichDoNot(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	production := envFixture(t, repo, "production", nil)
	staging := envFixture(t, repo, "staging", &production.ID)
	dev := envFixture(t, repo, "dev", &staging.ID)

	owner := configKeyFixture(t, repo, "owner", cdb.TypeString, true)
	noChildren := configKeyFixture(t, repo, "noChildren", cdb.TypeString, false)
	minReplicas := configKeyFixture(t, repo, "minReplicas", cdb.TypeInteger, true)
	maxReplicas := configKeyFixture(t, repo, "maxReplicas", cdb.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, cdb.NewIntConfigValue(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, maxReplicas.ID, 100))
	createConfigValue(t, repo, cdb.NewStringConfigValue(production.ID, noChildren.ID, "Nope"))

	expectedValues := []cdb.ConfigValue{
		*createConfigValue(t, repo, cdb.NewIntConfigValue(dev.ID, minReplicas.ID, 1)),
		*createInheritedConfigValue(t, repo, staging.Name, cdb.NewIntConfigValue(staging.ID, maxReplicas.ID, 50)),
		*createInheritedConfigValue(t, repo, production.Name, cdb.NewStringConfigValue(production.ID, owner.ID, "SRE")),
	}

	retrieved, err := repo.GetConfiguration(context.Background(), dev.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedValues, retrieved) {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", expectedValues, retrieved)
	}
}

func TestGetConfigurationShowsCanPropagateFalseKeysSetOnBaseEnvironment(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	production := envFixture(t, repo, "production", nil)
	staging := envFixture(t, repo, "staging", &production.ID)
	dev := envFixture(t, repo, "dev", &staging.ID)

	owner := configKeyFixture(t, repo, "owner", cdb.TypeString, true)
	noChildren := configKeyFixture(t, repo, "noChildren", cdb.TypeString, false)
	minReplicas := configKeyFixture(t, repo, "minReplicas", cdb.TypeInteger, true)
	maxReplicas := configKeyFixture(t, repo, "maxReplicas", cdb.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, cdb.NewIntConfigValue(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, maxReplicas.ID, 100))
	createConfigValue(t, repo, cdb.NewStringConfigValue(production.ID, noChildren.ID, "Nope"))

	expectedValues := []cdb.ConfigValue{
		*createConfigValue(t, repo, cdb.NewStringConfigValue(dev.ID, noChildren.ID, "Yes")),
		*createConfigValue(t, repo, cdb.NewIntConfigValue(dev.ID, minReplicas.ID, 1)),
		*createInheritedConfigValue(t, repo, staging.Name, cdb.NewIntConfigValue(staging.ID, maxReplicas.ID, 50)),
		*createInheritedConfigValue(t, repo, production.Name, cdb.NewStringConfigValue(production.ID, owner.ID, "SRE")),
	}

	retrieved, err := repo.GetConfiguration(context.Background(), dev.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedValues, retrieved) {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", expectedValues, retrieved)
	}
}

func TestGetConfigurationMarksInheritedValuesAsSuch(t *testing.T) {
	repo, tr := initTestDB(t)
	defer tr.Cleanup()

	production := envFixture(t, repo, "production", nil)
	staging := envFixture(t, repo, "staging", &production.ID)
	dev := envFixture(t, repo, "dev", &staging.ID)

	owner := configKeyFixture(t, repo, "owner", cdb.TypeString, true)
	minReplicas := configKeyFixture(t, repo, "minReplicas", cdb.TypeInteger, true)
	maxReplicas := configKeyFixture(t, repo, "maxReplicas", cdb.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, cdb.NewIntConfigValue(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, cdb.NewIntConfigValue(production.ID, maxReplicas.ID, 100))

	expectedValues := []cdb.ConfigValue{
		*createConfigValue(t, repo, cdb.NewIntConfigValue(dev.ID, minReplicas.ID, 1)),
		*createInheritedConfigValue(t, repo, staging.Name, cdb.NewIntConfigValue(staging.ID, maxReplicas.ID, 50)),
		*createInheritedConfigValue(t, repo, production.Name, cdb.NewStringConfigValue(production.ID, owner.ID, "SRE")),
	}

	retrieved, err := repo.GetConfiguration(context.Background(), dev.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedValues, retrieved) {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", expectedValues, retrieved)
	}
}
