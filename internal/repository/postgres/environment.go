package postgres

import (
	"context"
	_ "embed"
	"errors"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/internal/postgresutils"
	"github.com/jackc/pgx/v5"
)

//go:embed queries/environments/create_environment.sql
var createEnvironmentSql string

//go:embed queries/environments/get_environment_by_id.sql
var getEnvironmentByIDSql string

//go:embed queries/environments/get_environment_by_name.sql
var getEnvironmentByNameSql string

//go:embed queries/environments/list_environments.sql
var listEnvironmentsSql string

func (r *Repository) CreateEnvironment(ctx context.Context, env cdb.Environment) (cdb.Environment, error) {
	return postgresutils.GetOne[cdb.Environment](r.pool, ctx, createEnvironmentSql, env.Name, env.PromotesToID)
}

func (r *Repository) GetEnvironment(ctx context.Context, id int) (cdb.Environment, error) {
	env, err := postgresutils.GetOne[cdb.Environment](r.pool, ctx, getEnvironmentByIDSql, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return env, cdb.ErrEnvNotFound
	}

	return env, err
}

func (r *Repository) GetEnvironmentByName(ctx context.Context, name string) (cdb.Environment, error) {
	env, err := postgresutils.GetOne[cdb.Environment](r.pool, ctx, getEnvironmentByNameSql, name)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return env, cdb.ErrEnvNotFound
	}

	return env, err
}

func (r *Repository) ListEnvironments(ctx context.Context) ([]cdb.Environment, error) {
	environs, err := postgresutils.GetAll[cdb.Environment](r.pool, ctx, listEnvironmentsSql)
	return environs, err
}
