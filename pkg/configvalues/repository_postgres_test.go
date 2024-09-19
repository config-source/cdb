package configvalues

import (
	"context"
	_ "embed"
	"errors"
	"reflect"
	"testing"

	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/postgresutils"
	"github.com/config-source/cdb/pkg/services"
	"github.com/rs/zerolog"
)

func initTestDB(t *testing.T) (
	Repository,
	*environments.PostgresRepository,
	*configkeys.PostgresRepository,
	*services.PostgresRepository,
	*postgresutils.TestRepository,
) {
	t.Helper()

	tr, pool := postgresutils.InitTestDB(t)
	logger := zerolog.New(nil).Level(zerolog.Disabled)

	envRepo := environments.NewRepository(logger, pool)
	keyRepo := configkeys.NewRepository(logger, pool)
	svcRepo := services.NewRepository(logger, pool)
	repo := NewRepository(logger, pool, envRepo)

	return repo, envRepo, keyRepo, svcRepo, tr
}

func envFixture(
	t *testing.T,
	repo *environments.PostgresRepository,
	name string,
	promotesToID *int,
	serviceID int,
) environments.Environment {
	env, err := repo.CreateEnvironment(context.Background(), environments.Environment{
		Name:         name,
		PromotesToID: promotesToID,
		ServiceID:    serviceID,
	})
	if err != nil {
		t.Fatal(err)
	}

	return env
}

func svcFixture(t *testing.T, repo *services.PostgresRepository, name string) services.Service {
	svc, err := repo.CreateService(context.Background(), services.Service{
		Name: name,
	})
	if err != nil {
		t.Fatal(err)
	}

	return svc
}

func configKeyFixture(
	t *testing.T,
	repo *configkeys.PostgresRepository,
	svcID int,
	name string,
	valueType configkeys.ValueType,
	canPropagate bool,
) configkeys.ConfigKey {
	ck, err := repo.CreateConfigKey(
		context.Background(),
		configkeys.NewWithPropagation(
			svcID,
			name,
			valueType,
			canPropagate,
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	return ck
}

func TestCreateConfigValue(t *testing.T) {
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	env := envFixture(t, envRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeString, true)

	val := "test"
	cv, err := repo.CreateConfigValue(
		context.Background(),
		NewString(env.ID, key.ID, val),
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
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	env := envFixture(t, envRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	cv := createConfigValue(t, repo, NewString(
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
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	env := envFixture(t, envRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	cv := createConfigValue(t, repo, NewString(
		env.ID,
		key.ID,
		"test",
	))
	cv.ID = cv.ID + 1

	var err error
	_, err = repo.UpdateConfigurationValue(context.Background(), cv)
	expectedError := ErrNotFound
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected %s Got %s", expectedError, err)
	}
}

func TestUpdateConfigValueReturnsErrConfigKeyNotFound(t *testing.T) {
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	env := envFixture(t, envRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	cv := createConfigValue(t, repo, NewString(
		env.ID,
		key.ID,
		"test",
	))
	cv.ConfigKeyID = cv.ConfigKeyID + 1

	var err error
	_, err = repo.UpdateConfigurationValue(context.Background(), cv)
	expectedError := configkeys.ErrNotFound
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected %s Got %s", expectedError, err)
	}
}

func TestUpdateConfigValueReturnsErrEnvironmentNotFound(t *testing.T) {
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	env := envFixture(t, envRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	cv := createConfigValue(t, repo, NewString(
		env.ID,
		key.ID,
		"test",
	))
	cv.EnvironmentID = cv.EnvironmentID + 1

	var err error
	_, err = repo.UpdateConfigurationValue(context.Background(), cv)
	expectedError := environments.ErrNotFound
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected %s Got %s", expectedError, err)
	}
}

func TestCreateIntConfigValue(t *testing.T) {
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	env := envFixture(t, envRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeInteger, true)

	val := "test"
	cv, err := repo.CreateConfigValue(
		context.Background(),
		NewString(env.ID, key.ID, val),
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
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	env := envFixture(t, envRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	secondKey := configKeyFixture(t, keyRepo, svc.ID, "secondKey", configkeys.TypeString, true)

	_, err := repo.CreateConfigValue(context.Background(), NewString(
		env.ID,
		secondKey.ID,
		"test",
	))
	if err != nil {
		t.Fatal(err)
	}

	cv, err := repo.CreateConfigValue(context.Background(), NewString(
		env.ID,
		key.ID,
		"test",
	))
	if err != nil {
		t.Fatal(err)
	}

	retrieved, err := repo.GetConfigurationValue(context.Background(), env.ID, key.Name)
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
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	production := envFixture(t, envRepo, "production", nil, svc.ID)
	staging := envFixture(t, envRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, envRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	minReplicas := configKeyFixture(t, keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, NewInt(production.ID, maxReplicas.ID, 100))

	setDirectly := createConfigValue(t, repo, NewInt(dev.ID, minReplicas.ID, 1))
	stagingInherited := createInheritedConfigValue(t, repo, staging.Name, NewInt(staging.ID, maxReplicas.ID, 50))
	productionInherited := createInheritedConfigValue(t, repo, production.Name, NewString(production.ID, owner.ID, "SRE"))

	setDirectlyValue, err := repo.GetConfigurationValue(context.Background(), dev.ID, setDirectly.Name)
	if err != nil {
		t.Fatal(err)
	}

	if setDirectly.ID != setDirectlyValue.ID {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", setDirectly, setDirectlyValue)
	}

	stagingValue, err := repo.GetConfigurationValue(context.Background(), dev.ID, stagingInherited.Name)
	if err != nil {
		t.Fatal(err)
	}

	if stagingInherited.ID != stagingValue.ID {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", stagingInherited, stagingValue)
	}

	if !stagingValue.Inherited {
		t.Fatalf("Expected staging inherited value to be marked as such: %t\n", stagingValue.Inherited)
	}

	productionValue, err := repo.GetConfigurationValue(context.Background(), dev.ID, productionInherited.Name)
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
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	production := envFixture(t, envRepo, "production", nil, svc.ID)
	staging := envFixture(t, envRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, envRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	minReplicas := configKeyFixture(t, keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, NewInt(production.ID, maxReplicas.ID, 100))
	createConfigValue(t, repo, NewInt(dev.ID, minReplicas.ID, 1))
	createInheritedConfigValue(t, repo, staging.Name, NewInt(staging.ID, maxReplicas.ID, 50))
	createInheritedConfigValue(t, repo, production.Name, NewString(production.ID, owner.ID, "SRE"))

	_, err := repo.GetConfigurationValue(context.Background(), dev.ID, "notfound")
	if err != ErrNotFound {
		t.Fatalf("Expected: %s Got: %s", ErrNotFound, err)
	}
}

func TestGetConfigValueReturnsCorrectErrorForEnvNotFound(t *testing.T) {
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	production := envFixture(t, envRepo, "production", nil, svc.ID)
	staging := envFixture(t, envRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, envRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	minReplicas := configKeyFixture(t, keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, NewInt(production.ID, maxReplicas.ID, 100))
	createConfigValue(t, repo, NewInt(dev.ID, minReplicas.ID, 1))
	createInheritedConfigValue(t, repo, staging.Name, NewInt(staging.ID, maxReplicas.ID, 50))
	createInheritedConfigValue(t, repo, production.Name, NewString(production.ID, owner.ID, "SRE"))

	_, err := repo.GetConfigurationValue(context.Background(), 1000, "notfound")
	if !errors.Is(err, environments.ErrNotFound) {
		t.Fatalf("Expected: %s Got: %s", environments.ErrNotFound, err)
	}
}

func createConfigValue(t *testing.T, repo Repository, cv *ConfigValue) *ConfigValue {
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

func createInheritedConfigValue(t *testing.T, repo Repository, parentName string, cv *ConfigValue) *ConfigValue {
	created := createConfigValue(t, repo, cv)
	created.Inherited = true
	created.InheritedFrom = parentName
	return created
}

func TestGetConfiguration(t *testing.T) {
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	production := envFixture(t, envRepo, "production", nil, svc.ID)
	staging := envFixture(t, envRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, envRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	minReplicas := configKeyFixture(t, keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, NewInt(production.ID, maxReplicas.ID, 100))

	expectedValues := []ConfigValue{
		*createConfigValue(t, repo, NewInt(dev.ID, minReplicas.ID, 1)),
		*createInheritedConfigValue(t, repo, staging.Name, NewInt(staging.ID, maxReplicas.ID, 50)),
		*createInheritedConfigValue(t, repo, production.Name, NewString(production.ID, owner.ID, "SRE")),
	}

	retrieved, err := repo.GetConfiguration(context.Background(), dev.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedValues, retrieved) {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", expectedValues, retrieved)
	}
}

func TestGetConfigurationDoesntPropagateKeysWhichDoNot(t *testing.T) {
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	production := envFixture(t, envRepo, "production", nil, svc.ID)
	staging := envFixture(t, envRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, envRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	noChildren := configKeyFixture(t, keyRepo, svc.ID, "noChildren", configkeys.TypeString, false)
	minReplicas := configKeyFixture(t, keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, NewInt(production.ID, maxReplicas.ID, 100))
	createConfigValue(t, repo, NewString(production.ID, noChildren.ID, "Nope"))

	expectedValues := []ConfigValue{
		*createConfigValue(t, repo, NewInt(dev.ID, minReplicas.ID, 1)),
		*createInheritedConfigValue(t, repo, staging.Name, NewInt(staging.ID, maxReplicas.ID, 50)),
		*createInheritedConfigValue(t, repo, production.Name, NewString(production.ID, owner.ID, "SRE")),
	}

	retrieved, err := repo.GetConfiguration(context.Background(), dev.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedValues, retrieved) {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", expectedValues, retrieved)
	}
}

func TestGetConfigurationShowsCanPropagateFalseKeysSetOnBaseEnvironment(t *testing.T) {
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	production := envFixture(t, envRepo, "production", nil, svc.ID)
	staging := envFixture(t, envRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, envRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	noChildren := configKeyFixture(t, keyRepo, svc.ID, "noChildren", configkeys.TypeString, false)
	minReplicas := configKeyFixture(t, keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, NewInt(production.ID, maxReplicas.ID, 100))
	createConfigValue(t, repo, NewString(production.ID, noChildren.ID, "Nope"))

	expectedValues := []ConfigValue{
		*createConfigValue(t, repo, NewString(dev.ID, noChildren.ID, "Yes")),
		*createConfigValue(t, repo, NewInt(dev.ID, minReplicas.ID, 1)),
		*createInheritedConfigValue(t, repo, staging.Name, NewInt(staging.ID, maxReplicas.ID, 50)),
		*createInheritedConfigValue(t, repo, production.Name, NewString(production.ID, owner.ID, "SRE")),
	}

	retrieved, err := repo.GetConfiguration(context.Background(), dev.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedValues, retrieved) {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", expectedValues, retrieved)
	}
}

func TestGetConfigurationMarksInheritedValuesAsSuch(t *testing.T) {
	repo, envRepo, keyRepo, svcRepo, tr := initTestDB(t)
	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	production := envFixture(t, envRepo, "production", nil, svc.ID)
	staging := envFixture(t, envRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, envRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	minReplicas := configKeyFixture(t, keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, repo, NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, repo, NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, repo, NewInt(production.ID, maxReplicas.ID, 100))

	expectedValues := []ConfigValue{
		*createConfigValue(t, repo, NewInt(dev.ID, minReplicas.ID, 1)),
		*createInheritedConfigValue(t, repo, staging.Name, NewInt(staging.ID, maxReplicas.ID, 50)),
		*createInheritedConfigValue(t, repo, production.Name, NewString(production.ID, owner.ID, "SRE")),
	}

	retrieved, err := repo.GetConfiguration(context.Background(), dev.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedValues, retrieved) {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", expectedValues, retrieved)
	}
}
