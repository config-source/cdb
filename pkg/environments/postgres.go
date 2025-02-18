package environments

import (
	"context"
	_ "embed"
	"errors"

	"github.com/config-source/cdb/pkg/postgresutils"
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

//go:embed queries/update_environment.sql
var updateEnvironmentSql string

//go:embed queries/delete_environment.sql
var deleteEnvironmentSql string

func (r *Repository) CreateEnvironment(ctx context.Context, env Environment) (Environment, error) {
	return postgresutils.GetOneLax[Environment](
		r.pool,
		ctx,
		createEnvironmentSql,
		env.Name,
		env.PromotesToID,
		env.Sensitive,
		env.ServiceID,
	)
}

func (r *Repository) GetEnvironment(ctx context.Context, id int) (Environment, error) {
	env, err := postgresutils.GetOne[Environment](r.pool, ctx, getEnvironmentByIDSql, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return env, ErrNotFound
	}

	return env, err
}

func (r *Repository) GetEnvironmentByName(ctx context.Context, serviceName, name string) (Environment, error) {
	env, err := postgresutils.GetOne[Environment](r.pool, ctx, getEnvironmentByNameSql, serviceName, name)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return env, ErrNotFound
	}

	return env, err
}

func (r *Repository) ListEnvironments(ctx context.Context, includeSensitive bool) ([]Environment, error) {
	sql := listNonsensitiveEnvironmentsSql
	if includeSensitive {
		sql = listEnvironmentsSql
	}

	envs, err := postgresutils.GetAll[Environment](r.pool, ctx, sql)
	r.log.Debug().Int("envCount", len(envs)).Msg("retrieved environments")
	return envs, err
}

func (r *Repository) UpdateEnvironment(ctx context.Context, env Environment) (Environment, error) {
	return postgresutils.GetOneLax[Environment](
		r.pool,
		ctx,
		updateEnvironmentSql,
		env.ID,
		env.Name,
		env.PromotesToID,
		env.Sensitive,
	)
}

func (r *Repository) DeleteEnvironment(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, deleteEnvironmentSql, id)
	return err
}
