package webhooks_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/postgresutils"
	"github.com/config-source/cdb/pkg/services"
	"github.com/config-source/cdb/pkg/webhooks"
	"github.com/rs/zerolog"
)

func initTestDB(t *testing.T) (webhooks.DefinitionRepository, *environments.PostgresRepository, *services.PostgresRepository, *postgresutils.TestRepository) {
	t.Helper()

	tr, pool := postgresutils.InitTestDB(t)
	repo := webhooks.NewRepository(
		zerolog.New(nil).Level(zerolog.Disabled),
		pool,
	)
	envRepo := environments.NewRepository(
		zerolog.New(nil).Level(zerolog.Disabled),
		pool,
	)
	svcRepo := services.NewRepository(
		zerolog.New(nil).Level(zerolog.Disabled),
		pool,
	)

	return repo, envRepo, svcRepo, tr
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

func TestCreateWebhookDefinition(t *testing.T) {
	repo, _, _, tr := initTestDB(t)

	defer tr.Cleanup()

	wh, err := repo.CreateWebhookDefinition(context.Background(), webhooks.Definition{
		Template: "Smile",
		URL:      "http://localhost:8081/webhooks",
	})

	if err != nil {
		t.Fatal(err)
	}

	if wh.ID == 0 {
		t.Fatalf("Expected ID to be set got: %d", wh.ID)
	}

	if wh.Template != "Smile" {
		t.Fatalf("Expected Template to be Smile got: %s", wh.Template)
	}

	if wh.URL != "http://localhost:8081/webhooks" {
		t.Fatalf("Expected URL to be http://localhost:8081/webhooks got: %s", wh.URL)
	}

	if wh.AuthzHeader != "" {
		t.Fatalf("Expected AuthzHeader to be empty got: %s", wh.AuthzHeader)
	}
}

func TestCreateWebhookDefinitionWithAuthzHeader(t *testing.T) {
	repo, _, _, tr := initTestDB(t)

	defer tr.Cleanup()

	wh, err := repo.CreateWebhookDefinition(context.Background(), webhooks.Definition{
		Template:    "Smile",
		URL:         "http://localhost:8081/webhooks",
		AuthzHeader: "letmein",
	})

	if err != nil {
		t.Fatal(err)
	}

	if wh.ID == 0 {
		t.Fatalf("Expected ID to be set got: %d", wh.ID)
	}

	if wh.Template != "Smile" {
		t.Fatalf("Expected Template to be Smile got: %s", wh.Template)
	}

	if wh.URL != "http://localhost:8081/webhooks" {
		t.Fatalf("Expected URL to be http://localhost:8081/webhooks got: %s", wh.URL)
	}

	if wh.AuthzHeader != "letmein" {
		t.Fatalf("Expected AuthzHeader to be letmein got: %s", wh.AuthzHeader)
	}
}

func TestDeleteWebhookDefinition(t *testing.T) {
	repo, _, _, tr := initTestDB(t)

	defer tr.Cleanup()

	wh, err := repo.CreateWebhookDefinition(context.Background(), webhooks.Definition{
		Template:    "Smile",
		URL:         "http://localhost:8081/webhooks",
		AuthzHeader: "letmein",
	})

	if err != nil {
		t.Fatal(err)
	}

	err = repo.DeleteWebhookDefinition(context.Background(), wh.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.GetWebhookDefinition(context.Background(), wh.ID)

	if !errors.Is(err, webhooks.ErrNotFound) {
		t.Fatalf("Expected %s, got: %s", webhooks.ErrNotFound, err)
	}
}

func TestGetWebhookDefinition(t *testing.T) {
	repo, _, _, tr := initTestDB(t)

	defer tr.Cleanup()

	wh, err := repo.CreateWebhookDefinition(context.Background(), webhooks.Definition{
		Template:    "Smile",
		URL:         "http://localhost:8081/webhooks",
		AuthzHeader: "letmein",
	})

	if err != nil {
		t.Fatal(err)
	}

	wh1, err := repo.GetWebhookDefinition(context.Background(), wh.ID)
	if errors.Is(err, webhooks.ErrNotFound) {
		t.Fatalf("Expected to find webhook definition, got: %s", webhooks.ErrNotFound)
	} else if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(wh, wh1) {
		t.Fatalf("Got wrong webhoook definition, expected %+v, got %+v", wh, wh1)
	}
}

func webhookDefinitionFixture(t *testing.T, repo webhooks.DefinitionRepository, wh webhooks.Definition, envIDs ...int) webhooks.Definition {
	t.Helper()
	wh, err := repo.CreateWebhookDefinition(context.Background(), wh)

	if err != nil {
		t.Fatal(err)
	}

	err = repo.AssociateWebhookToEnvironments(context.Background(), wh.ID, envIDs)
	if err != nil {
		t.Fatal(err)
	}

	return wh
}

func TestGetWebhookDefinitions(t *testing.T) {
	repo, envRepo, svcRepo, tr := initTestDB(t)

	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")
	env1 := envFixture(t, envRepo, "env1", nil, svc.ID)
	env2 := envFixture(t, envRepo, "env2", nil, svc.ID)

	webhookDefinitionFixture(
		t,
		repo,
		webhooks.Definition{
			Template:    "Frown",
			URL:         "http://localhost:8081/webhooks",
			AuthzHeader: "letmein",
		},
		env2.ID,
	)

	expectedWhs := []webhooks.Definition{
		webhookDefinitionFixture(t, repo, webhooks.Definition{
			Template:    "Smile",
			URL:         "http://localhost:8081/webhooks",
			AuthzHeader: "letmein",
		}, env1.ID),
		webhookDefinitionFixture(t, repo, webhooks.Definition{
			Template:    "Smile2",
			URL:         "http://localhost:8082/webhooks",
			AuthzHeader: "letmein2",
		}, env1.ID),
	}

	retreivedWhs, err := repo.GetWebhookDefinitionsForEnvironment(context.Background(), env1.ID) // TODO make this actually work
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedWhs, retreivedWhs) {
		t.Fatalf("Got wrong webhoook definition, expected %+v, got %+v", expectedWhs, retreivedWhs)
	}
}

func TestAssociateWebhookWithSingleEnvironment(t *testing.T) {
	repo, envRepo, svcRepo, tr := initTestDB(t)

	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")

	env1 := envFixture(t, envRepo, "env1", nil, svc.ID)
	wh, err := repo.CreateWebhookDefinition(context.Background(), webhooks.Definition{
		Template:    "Smile",
		URL:         "http://localhost:8081/webhooks",
		AuthzHeader: "letmein",
	})

	if err != nil {
		t.Fatal(err)
	}

	err = repo.AssociateWebhookToEnvironments(context.Background(), wh.ID, []int{env1.ID})
	if err != nil {
		t.Fatal(err)
	}

	var count int
	err = tr.Pool.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM webhook_definitions_to_environments WHERE webhook_definition_id = $1 AND environment_id = $2",
		wh.ID,
		env1.ID,
	).Scan(&count)

	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatalf("Expected 1, got %d", count)
	}
}

func TestAssociateWebhookWithMultipleEnvironments(t *testing.T) {
	repo, envRepo, svcRepo, tr := initTestDB(t)

	defer tr.Cleanup()

	svc := svcFixture(t, svcRepo, "svc1")

	env1 := envFixture(t, envRepo, "env1", nil, svc.ID)
	env2 := envFixture(t, envRepo, "env2", nil, svc.ID)
	env3 := envFixture(t, envRepo, "env3", nil, svc.ID)

	wh, err := repo.CreateWebhookDefinition(context.Background(), webhooks.Definition{
		Template:    "Smile",
		URL:         "http://localhost:8081/webhooks",
		AuthzHeader: "letmein",
	})

	if err != nil {
		t.Fatal(err)
	}

	err = repo.AssociateWebhookToEnvironments(context.Background(), wh.ID, []int{env1.ID, env2.ID, env3.ID})
	if err != nil {
		t.Fatal(err)
	}

	var count int
	err = tr.Pool.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM webhook_definitions_to_environments WHERE webhook_definition_id = $1",
		wh.ID,
	).Scan(&count)

	if err != nil {
		t.Fatal(err)
	}

	if count != 3 {
		t.Fatalf("Expected 3, got %d", count)
	}
}
