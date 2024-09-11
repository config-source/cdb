package configkeys

import (
	"context"
	_ "embed"
	"errors"

	"github.com/config-source/cdb/postgresutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func NewRepository(log zerolog.Logger, pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		log:  log,
		pool: pool,
	}
}

//go:embed queries/create_config_key.sql
var createConfigKeySql string

//go:embed queries/get_config_key_by_id.sql
var getConfigKeyByIDSql string

//go:embed queries/get_config_key_by_name.sql
var getConfigKeyByNameSql string

//go:embed queries/get_all_config_keys.sql
var getAllConfigKeys string

func (r *PostgresRepository) CreateConfigKey(ctx context.Context, ck ConfigKey) (ConfigKey, error) {
	var canPropagate bool
	if ck.CanPropagate == nil {
		canPropagate = true
	} else {
		canPropagate = *ck.CanPropagate
	}

	return postgresutils.GetOne[ConfigKey](r.pool, ctx, createConfigKeySql, ck.Name, ck.ValueType, canPropagate)
}

func (r *PostgresRepository) GetConfigKey(ctx context.Context, id int) (ConfigKey, error) {
	key, err := postgresutils.GetOne[ConfigKey](r.pool, ctx, getConfigKeyByIDSql, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return key, ErrNotFound
	}

	return key, err
}

func (r *PostgresRepository) GetConfigKeyByName(ctx context.Context, name string) (ConfigKey, error) {
	key, err := postgresutils.GetOne[ConfigKey](r.pool, ctx, getConfigKeyByNameSql, name)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return key, ErrNotFound
	}

	return key, err

}

func (r *PostgresRepository) ListConfigKeys(ctx context.Context) ([]ConfigKey, error) {
	return postgresutils.GetAll[ConfigKey](r.pool, ctx, getAllConfigKeys)
}