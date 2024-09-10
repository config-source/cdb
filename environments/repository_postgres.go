package environments

import (
	"context"
	_ "embed"
	"errors"

	"github.com/config-source/cdb/postgresutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Repository struct {
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func NewRepository(log zerolog.Logger, pool *pgxpool.Pool) *Repository {
	return &Repository{
		log:  log,
		pool: pool,
	}
}

//go:embed queries/create_environment.sql
var createEnvironmentSql string

//go:embed queries/get_environment_by_id.sql
var getEnvironmentByIDSql string

//go:embed queries/get_environment_by_name.sql
var getEnvironmentByNameSql string

//go:embed queries/list_environments.sql
var listEnvironmentsSql string

//go:embed queries/list_nonsensitive_environments.sql
var listNonsensitiveEnvironmentsSql string

func (r *Repository) CreateEnvironment(ctx context.Context, env Environment) (Environment, error) {
	return postgresutils.GetOne[Environment](r.pool, ctx, createEnvironmentSql, env.Name, env.PromotesToID, env.Sensitive)
}

func (r *Repository) GetEnvironment(ctx context.Context, id int) (Environment, error) {
	env, err := postgresutils.GetOne[Environment](r.pool, ctx, getEnvironmentByIDSql, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return env, ErrEnvNotFound
	}

	return env, err
}

func (r *Repository) GetEnvironmentByName(ctx context.Context, name string) (Environment, error) {
	env, err := postgresutils.GetOne[Environment](r.pool, ctx, getEnvironmentByNameSql, name)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return env, ErrEnvNotFound
	}

	return env, err
}

func (r *Repository) ListEnvironments(ctx context.Context, includeSensitive bool) ([]Environment, error) {
	sql := listNonsensitiveEnvironmentsSql
	if includeSensitive {
		sql = listEnvironmentsSql
	}

	return postgresutils.GetAll[Environment](r.pool, ctx, sql)
}
