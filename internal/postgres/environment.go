package postgres

import (
	"context"
	_ "embed"

	"github.com/config-source/cdb"
)

//go:embed queries/environments/create_environment.sql
var createEnvironmentSql string

//go:embed queries/environments/get_environment_by_id.sql
var getEnvironmentByIDSql string

//go:embed queries/environments/get_environment_by_name.sql
var getEnvironmentByNameSql string

func (r *Repository) CreateEnvironment(ctx context.Context, env cdb.Environment) (cdb.Environment, error) {
	return getOne[cdb.Environment](r, ctx, createEnvironmentSql, env.Name, env.PromotesToID)
}

func (r *Repository) GetEnvironment(ctx context.Context, id int) (cdb.Environment, error) {
	env, err := getOne[cdb.Environment](r, ctx, getEnvironmentByIDSql, id)
	if err != nil && err.Error() == "no rows in result set" {
		return env, cdb.ErrEnvNotFound
	}

	return env, err
}

func (r *Repository) GetEnvironmentByName(ctx context.Context, name string) (cdb.Environment, error) {
	env, err := getOne[cdb.Environment](r, ctx, getEnvironmentByNameSql, name)
	if err != nil && err.Error() == "no rows in result set" {
		return env, cdb.ErrEnvNotFound
	}

	return env, err
}
