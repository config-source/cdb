package postgres

import (
	"context"
	_ "embed"
	"errors"

	"github.com/config-source/cdb"
	"github.com/jackc/pgx/v5"
)

//go:embed queries/webhookDefinitions/create_webhook_definition.sql
var createWebhookDefinitionSql string

//go:embed queries/webhookDefinitions/delete_webhook_definition.sql
var deleteWebhookDefinitionSql string

//go:embed queries/webhookDefinitions/get_webhook_definition_by_id.sql
var getWebhookDefinitionByIDSql string

//go:embed queries/webhookDefinitions/get_all_webhook_definitions_for_environment.sql
var getAllWebhookDefinitionsForEnvironmentSql string

func (r *Repository) CreateWebhookDefinition(ctx context.Context, wh cdb.WebhookDefinition) (cdb.WebhookDefinition, error) {
	return getOne[cdb.WebhookDefinition](r, ctx, createWebhookDefinitionSql, wh.Template, wh.URL, wh.AuthzHeader)
}

func (r *Repository) DeleteWebhookDefinition(ctx context.Context, webhookID int) error {
	_, err := r.pool.Exec(ctx, deleteWebhookDefinitionSql, webhookID)
	return err
}

func (r *Repository) GetWebhookDefinition(ctx context.Context, id int) (cdb.WebhookDefinition, error) {
	key, err := getOne[cdb.WebhookDefinition](r, ctx, getWebhookDefinitionByIDSql, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return key, cdb.ErrWebhookDefinitionNotFound
	}

	return key, err
}

func (r *Repository) GetWebhookDefinitions(ctx context.Context, envID int) ([]cdb.WebhookDefinition, error) {
	return getAll[cdb.WebhookDefinition](r, ctx, getAllWebhookDefinitionsForEnvironmentSql)
}
