package webhooks

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/config-source/cdb/pkg/postgresutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func NewRepository(log zerolog.Logger, pool *pgxpool.Pool) DefinitionRepository {
	return &PostgresRepository{
		log:  log,
		pool: pool,
	}
}

//go:embed queries/create_webhook_definition.sql
var createWebhookDefinitionSql string

//go:embed queries/delete_webhook_definition.sql
var deleteWebhookDefinitionSql string

//go:embed queries/get_webhook_definition_by_id.sql
var getWebhookDefinitionByIDSql string

//go:embed queries/get_all_webhook_definitions_for_environment.sql
var getAllWebhookDefinitionsForEnvironmentSql string

func (r *PostgresRepository) CreateWebhookDefinition(ctx context.Context, wh Definition) (Definition, error) {
	return postgresutils.GetOneLax[Definition](r.pool, ctx, createWebhookDefinitionSql, wh.Template, wh.URL, wh.AuthzHeader)
}

func (r *PostgresRepository) AssociateWebhookToEnvironments(ctx context.Context, webhookID int, environmentIDs []int) error {
	rows := make([][]any, len(environmentIDs))
	for i, environmentID := range environmentIDs {
		rows[i] = []any{webhookID, environmentID}
	}

	count, err := r.pool.CopyFrom(
		ctx,
		pgx.Identifier{"webhook_definitions_to_environments"},
		[]string{"webhook_definition_id", "environment_id"},
		pgx.CopyFromRows(rows),
	)

	if err == nil && int64(len(environmentIDs)) != count {
		return fmt.Errorf("expected to create %d associations, instead created %d", len(environmentIDs), count)
	}

	return err
}

func (r *PostgresRepository) DeleteWebhookDefinition(ctx context.Context, webhookID int) error {
	_, err := r.pool.Exec(ctx, deleteWebhookDefinitionSql, webhookID)
	return err
}

func (r *PostgresRepository) GetWebhookDefinition(ctx context.Context, id int) (Definition, error) {
	key, err := postgresutils.GetOne[Definition](r.pool, ctx, getWebhookDefinitionByIDSql, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return key, ErrNotFound
	}

	return key, err
}

func (r *PostgresRepository) GetWebhookDefinitionsForEnvironment(ctx context.Context, envID int) ([]Definition, error) {
	return postgresutils.GetAll[Definition](r.pool, ctx, getAllWebhookDefinitionsForEnvironmentSql, envID)
}
