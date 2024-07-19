package webhooks

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("webhook defintion not found")
)

type Definition struct {
	ID          int       `db:"id"`
	Template    string    `db:"template"`
	URL         string    `db:"url"`
	AuthzHeader string    `db:"authz_header"`
	CreatedAt   time.Time `db:"created_at"`
}

type DefinitionRepository interface {
	GetWebhookDefinition(ctx context.Context, id int) (Definition, error)
	GetWebhookDefinitionsForEnvironment(ctx context.Context, envID int) ([]Definition, error)
	CreateWebhookDefinition(context.Context, Definition) (Definition, error)
	AssociateWebhookToEnvironments(ctx context.Context, webhookID int, environmentIDs []int) error
	// UpdateWebhookDefinition(context.Context, WebhookDefinition) (WebhookDefinition, error)
	DeleteWebhookDefinition(ctx context.Context, webhookID int) error
}
