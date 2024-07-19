package postgres_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/internal/postgres"
)

func TestCreateWebhookDefinition(t *testing.T) {
	repo, tr := initTestDB(t)

	defer tr.Cleanup()

	wh, err := repo.CreateWebhookDefinition(context.Background(), cdb.WebhookDefinition{
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
	repo, tr := initTestDB(t)

	defer tr.Cleanup()

	wh, err := repo.CreateWebhookDefinition(context.Background(), cdb.WebhookDefinition{
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
	repo, tr := initTestDB(t)

	defer tr.Cleanup()

	wh, err := repo.CreateWebhookDefinition(context.Background(), cdb.WebhookDefinition{
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

	if !errors.Is(err, cdb.ErrWebhookDefinitionNotFound) {
		t.Fatalf("Expected %s, got: %s", cdb.ErrWebhookDefinitionNotFound, err)
	}
}

func TestGetWebhookDefinition(t *testing.T) {
	repo, tr := initTestDB(t)

	defer tr.Cleanup()

	wh, err := repo.CreateWebhookDefinition(context.Background(), cdb.WebhookDefinition{
		Template:    "Smile",
		URL:         "http://localhost:8081/webhooks",
		AuthzHeader: "letmein",
	})

	if err != nil {
		t.Fatal(err)
	}

	wh1, err := repo.GetWebhookDefinition(context.Background(), wh.ID)
	if errors.Is(err, cdb.ErrWebhookDefinitionNotFound) {
		t.Fatalf("Expected to find webhook definition, got: %s", cdb.ErrWebhookDefinitionNotFound)
	} else if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(wh, wh1) {
		t.Fatalf("Got wrong webhoook definition, expected %+v, got %+v", wh, wh1)
	}
}

func createWebhookDefinition(t *testing.T, repo *postgres.Repository, wh cdb.WebhookDefinition) cdb.WebhookDefinition {
	wh, err := repo.CreateWebhookDefinition(context.Background(), wh)

	if err != nil {
		t.Fatal(err)
	}

	return wh
}

func TestGetWebhookDefinitions(t *testing.T) {
	repo, tr := initTestDB(t)

	defer tr.Cleanup()

	expectedWhs := []cdb.WebhookDefinition{
		createWebhookDefinition(t, repo, cdb.WebhookDefinition{
			Template:    "Smile",
			URL:         "http://localhost:8081/webhooks",
			AuthzHeader: "letmein",
		}),
		createWebhookDefinition(t, repo, cdb.WebhookDefinition{
			Template:    "Smile2",
			URL:         "http://localhost:8082/webhooks",
			AuthzHeader: "letmein2",
		}),
	}

	retreivedWhs, err := repo.GetWebhookDefinitions(context.Background(), 1) // TODO make this actually work
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedWhs, retreivedWhs) {
		t.Fatalf("Got wrong webhoook definition, expected %+v, got %+v", expectedWhs, retreivedWhs)
	}
}
