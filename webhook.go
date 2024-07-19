package cdb

import (
	"context"
	"errors"
	"time"
)

var (
	ErrWebhookDefinitionNotFound = errors.New("webhook defintion not found")
)

type WebhookDefinition struct {
	ID          int       `db:"id"`
	Template    string    `db:"template"`
	URL         string    `db:"url"`
	AuthzHeader string    `db:"authz_header"`
	CreatedAt   time.Time `db:"created_at"`
}

type WebhookDefinitionRepository interface {
	GetWebhookDefinitions(ctx context.Context, envID int) ([]WebhookDefinition, error)
	CreateWebhookDefinition(context.Context, WebhookDefinition) (WebhookDefinition, error)
	// UpdateWebhookDefinition(context.Context, WebhookDefinition) (WebhookDefinition, error)
	DeleteWebhookDefinition(ctx context.Context, webhookID int) error
}
