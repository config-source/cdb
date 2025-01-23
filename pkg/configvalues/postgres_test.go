package configvalues_test

import (
	"context"
	_ "embed"
	"errors"
	"reflect"
	"testing"

	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/configvalues"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/postgresutils"
	"github.com/config-source/cdb/pkg/services"
	"github.com/rs/zerolog"
)

type TestContext struct {
	valueRepo       *configvalues.Repository
	environmentRepo *environments.Repository
	keyRepo         *configkeys.Repository
	serviceRepo     *services.Repository
}

func initTestDB(t *testing.T) TestContext {
	t.Helper()

	pool := postgresutils.InitTestDB(t)
	logger := zerolog.New(nil).Level(zerolog.Disabled)

	envRepo := environments.NewRepository(logger, pool)
	keyRepo := configkeys.NewRepository(logger, pool)
	svcRepo := services.NewRepository(logger, pool)
	repo := configvalues.NewRepository(logger, pool, envRepo)

	return TestContext{
		valueRepo:       repo,
		environmentRepo: envRepo,
		keyRepo:         keyRepo,
		serviceRepo:     svcRepo,
	}
}

func envFixture(
	t *testing.T,
	repo *environments.Repository,
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

func svcFixture(t *testing.T, repo *services.Repository, name string) services.Service {
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
	repo *configkeys.Repository,
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
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	env := envFixture(t, tc.environmentRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeString, true)

	val := "test"
	cv, err := tc.valueRepo.CreateConfigValue(
		context.Background(),
		configvalues.NewString(env.ID, key.ID, val),
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
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	env := envFixture(t, tc.environmentRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	cv := createConfigValue(t, tc.valueRepo, configvalues.NewString(
		env.ID,
		key.ID,
		"test",
	))
	cv.SetStrValue("updated")

	var err error
	cv, err = tc.valueRepo.UpdateConfigurationValue(context.Background(), cv)
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
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	env := envFixture(t, tc.environmentRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	cv := createConfigValue(t, tc.valueRepo, configvalues.NewString(
		env.ID,
		key.ID,
		"test",
	))
	cv.ID = cv.ID + 1

	var err error
	_, err = tc.valueRepo.UpdateConfigurationValue(context.Background(), cv)
	expectedError := configvalues.ErrNotFound
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected %s Got %s", expectedError, err)
	}
}

func TestUpdateConfigValueReturnsErrConfigKeyNotFound(t *testing.T) {
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	env := envFixture(t, tc.environmentRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	cv := createConfigValue(t, tc.valueRepo, configvalues.NewString(
		env.ID,
		key.ID,
		"test",
	))
	cv.ConfigKeyID = cv.ConfigKeyID + 1

	var err error
	_, err = tc.valueRepo.UpdateConfigurationValue(context.Background(), cv)
	expectedError := configkeys.ErrNotFound
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected %s Got %s", expectedError, err)
	}
}

func TestUpdateConfigValueReturnsErrEnvironmentNotFound(t *testing.T) {
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	env := envFixture(t, tc.environmentRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	cv := createConfigValue(t, tc.valueRepo, configvalues.NewString(
		env.ID,
		key.ID,
		"test",
	))
	cv.EnvironmentID = cv.EnvironmentID + 1

	var err error
	_, err = tc.valueRepo.UpdateConfigurationValue(context.Background(), cv)
	expectedError := environments.ErrNotFound
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected %s Got %s", expectedError, err)
	}
}

func TestCreateIntConfigValue(t *testing.T) {
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	env := envFixture(t, tc.environmentRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeInteger, true)

	val := "test"
	cv, err := tc.valueRepo.CreateConfigValue(
		context.Background(),
		configvalues.NewString(env.ID, key.ID, val),
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
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	env := envFixture(t, tc.environmentRepo, "cdb", nil, svc.ID)
	key := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	secondKey := configKeyFixture(t, tc.keyRepo, svc.ID, "secondKey", configkeys.TypeString, true)

	_, err := tc.valueRepo.CreateConfigValue(context.Background(), configvalues.NewString(
		env.ID,
		secondKey.ID,
		"test",
	))
	if err != nil {
		t.Fatal(err)
	}

	cv, err := tc.valueRepo.CreateConfigValue(context.Background(), configvalues.NewString(
		env.ID,
		key.ID,
		"test",
	))
	if err != nil {
		t.Fatal(err)
	}

	retrieved, err := tc.valueRepo.GetConfigurationValue(context.Background(), env.ID, key.Name)
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
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	production := envFixture(t, tc.environmentRepo, "production", nil, svc.ID)
	staging := envFixture(t, tc.environmentRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, tc.environmentRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	minReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, maxReplicas.ID, 100))

	setDirectly := createConfigValue(t, tc.valueRepo, configvalues.NewInt(dev.ID, minReplicas.ID, 1))
	stagingInherited := createInheritedConfigValue(t, tc.valueRepo, staging.Name, configvalues.NewInt(staging.ID, maxReplicas.ID, 50))
	productionInherited := createInheritedConfigValue(t, tc.valueRepo, production.Name, configvalues.NewString(production.ID, owner.ID, "SRE"))

	setDirectlyValue, err := tc.valueRepo.GetConfigurationValue(context.Background(), dev.ID, setDirectly.Name)
	if err != nil {
		t.Fatal(err)
	}

	if setDirectly.ID != setDirectlyValue.ID {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", setDirectly, setDirectlyValue)
	}

	stagingValue, err := tc.valueRepo.GetConfigurationValue(context.Background(), dev.ID, stagingInherited.Name)
	if err != nil {
		t.Fatal(err)
	}

	if stagingInherited.ID != stagingValue.ID {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", stagingInherited, stagingValue)
	}

	if !stagingValue.Inherited {
		t.Fatalf("Expected staging inherited value to be marked as such: %t\n", stagingValue.Inherited)
	}

	productionValue, err := tc.valueRepo.GetConfigurationValue(context.Background(), dev.ID, productionInherited.Name)
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
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	production := envFixture(t, tc.environmentRepo, "production", nil, svc.ID)
	staging := envFixture(t, tc.environmentRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, tc.environmentRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	minReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, maxReplicas.ID, 100))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(dev.ID, minReplicas.ID, 1))
	createInheritedConfigValue(t, tc.valueRepo, staging.Name, configvalues.NewInt(staging.ID, maxReplicas.ID, 50))
	createInheritedConfigValue(t, tc.valueRepo, production.Name, configvalues.NewString(production.ID, owner.ID, "SRE"))

	_, err := tc.valueRepo.GetConfigurationValue(context.Background(), dev.ID, "notfound")
	if err != configvalues.ErrNotFound {
		t.Fatalf("Expected: %s Got: %s", configvalues.ErrNotFound, err)
	}
}

func TestGetConfigValueReturnsCorrectErrorForEnvNotFound(t *testing.T) {
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	production := envFixture(t, tc.environmentRepo, "production", nil, svc.ID)
	staging := envFixture(t, tc.environmentRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, tc.environmentRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	minReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, maxReplicas.ID, 100))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(dev.ID, minReplicas.ID, 1))
	createInheritedConfigValue(t, tc.valueRepo, staging.Name, configvalues.NewInt(staging.ID, maxReplicas.ID, 50))
	createInheritedConfigValue(t, tc.valueRepo, production.Name, configvalues.NewString(production.ID, owner.ID, "SRE"))

	_, err := tc.valueRepo.GetConfigurationValue(context.Background(), 1000, "notfound")
	if !errors.Is(err, environments.ErrNotFound) {
		t.Fatalf("Expected: %s Got: %s", environments.ErrNotFound, err)
	}
}

func createConfigValue(t *testing.T, repo *configvalues.Repository, cv *configvalues.ConfigValue) *configvalues.ConfigValue {
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

func createInheritedConfigValue(t *testing.T, repo *configvalues.Repository, parentName string, cv *configvalues.ConfigValue) *configvalues.ConfigValue {
	created := createConfigValue(t, repo, cv)
	created.Inherited = true
	created.InheritedFrom = parentName
	return created
}

func TestGetConfiguration(t *testing.T) {
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	production := envFixture(t, tc.environmentRepo, "production", nil, svc.ID)
	staging := envFixture(t, tc.environmentRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, tc.environmentRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	minReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, maxReplicas.ID, 100))

	expectedValues := []configvalues.ConfigValue{
		*createConfigValue(t, tc.valueRepo, configvalues.NewInt(dev.ID, minReplicas.ID, 1)),
		*createInheritedConfigValue(t, tc.valueRepo, staging.Name, configvalues.NewInt(staging.ID, maxReplicas.ID, 50)),
		*createInheritedConfigValue(t, tc.valueRepo, production.Name, configvalues.NewString(production.ID, owner.ID, "SRE")),
	}

	retrieved, err := tc.valueRepo.GetConfiguration(context.Background(), dev.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedValues, retrieved) {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", expectedValues, retrieved)
	}
}

func TestGetConfigurationDoesntPropagateKeysWhichDoNot(t *testing.T) {
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	production := envFixture(t, tc.environmentRepo, "production", nil, svc.ID)
	staging := envFixture(t, tc.environmentRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, tc.environmentRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	noChildren := configKeyFixture(t, tc.keyRepo, svc.ID, "noChildren", configkeys.TypeString, false)
	minReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, maxReplicas.ID, 100))
	createConfigValue(t, tc.valueRepo, configvalues.NewString(production.ID, noChildren.ID, "Nope"))

	expectedValues := []configvalues.ConfigValue{
		*createConfigValue(t, tc.valueRepo, configvalues.NewInt(dev.ID, minReplicas.ID, 1)),
		*createInheritedConfigValue(t, tc.valueRepo, staging.Name, configvalues.NewInt(staging.ID, maxReplicas.ID, 50)),
		*createInheritedConfigValue(t, tc.valueRepo, production.Name, configvalues.NewString(production.ID, owner.ID, "SRE")),
	}

	retrieved, err := tc.valueRepo.GetConfiguration(context.Background(), dev.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedValues, retrieved) {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", expectedValues, retrieved)
	}
}

func TestGetConfigurationShowsCanPropagateFalseKeysSetOnBaseEnvironment(t *testing.T) {
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	production := envFixture(t, tc.environmentRepo, "production", nil, svc.ID)
	staging := envFixture(t, tc.environmentRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, tc.environmentRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	noChildren := configKeyFixture(t, tc.keyRepo, svc.ID, "noChildren", configkeys.TypeString, false)
	minReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, maxReplicas.ID, 100))
	createConfigValue(t, tc.valueRepo, configvalues.NewString(production.ID, noChildren.ID, "Nope"))

	expectedValues := []configvalues.ConfigValue{
		*createConfigValue(t, tc.valueRepo, configvalues.NewString(dev.ID, noChildren.ID, "Yes")),
		*createConfigValue(t, tc.valueRepo, configvalues.NewInt(dev.ID, minReplicas.ID, 1)),
		*createInheritedConfigValue(t, tc.valueRepo, staging.Name, configvalues.NewInt(staging.ID, maxReplicas.ID, 50)),
		*createInheritedConfigValue(t, tc.valueRepo, production.Name, configvalues.NewString(production.ID, owner.ID, "SRE")),
	}

	retrieved, err := tc.valueRepo.GetConfiguration(context.Background(), dev.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedValues, retrieved) {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", expectedValues, retrieved)
	}
}

func TestGetConfigurationMarksInheritedValuesAsSuch(t *testing.T) {
	tc := initTestDB(t)

	svc := svcFixture(t, tc.serviceRepo, "svc1")
	production := envFixture(t, tc.environmentRepo, "production", nil, svc.ID)
	staging := envFixture(t, tc.environmentRepo, "staging", &production.ID, svc.ID)
	dev := envFixture(t, tc.environmentRepo, "dev", &staging.ID, svc.ID)

	owner := configKeyFixture(t, tc.keyRepo, svc.ID, "owner", configkeys.TypeString, true)
	minReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "minReplicas", configkeys.TypeInteger, true)
	maxReplicas := configKeyFixture(t, tc.keyRepo, svc.ID, "maxReplicas", configkeys.TypeInteger, true)

	// Throw in duplicate settings higher in the parent tree to ensure
	// inheritance overrides these values.
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(staging.ID, minReplicas.ID, 5))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, minReplicas.ID, 10))
	createConfigValue(t, tc.valueRepo, configvalues.NewInt(production.ID, maxReplicas.ID, 100))

	expectedValues := []configvalues.ConfigValue{
		*createConfigValue(t, tc.valueRepo, configvalues.NewInt(dev.ID, minReplicas.ID, 1)),
		*createInheritedConfigValue(t, tc.valueRepo, staging.Name, configvalues.NewInt(staging.ID, maxReplicas.ID, 50)),
		*createInheritedConfigValue(t, tc.valueRepo, production.Name, configvalues.NewString(production.ID, owner.ID, "SRE")),
	}

	retrieved, err := tc.valueRepo.GetConfiguration(context.Background(), dev.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedValues, retrieved) {
		t.Fatalf("\n\tExpected\n\t\t%+v\n\tGot\n\t\t%+v", expectedValues, retrieved)
	}
}
