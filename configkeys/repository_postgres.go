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

//go:embed queries/get_all_config_keys_by_service.sql
var getAllConfigKeysByService string

func (r *PostgresRepository) CreateConfigKey(ctx context.Context, ck ConfigKey) (ConfigKey, error) {
	var canPropagate bool
	if ck.CanPropagate == nil {
		canPropagate = true
	} else {
		canPropagate = *ck.CanPropagate
	}

	return postgresutils.GetOneLax[ConfigKey](
		r.pool,
		ctx,
		createConfigKeySql,
		ck.Name,
		ck.ValueType,
		canPropagate,
		ck.ServiceID,
	)
}

func (r *PostgresRepository) GetConfigKey(ctx context.Context, serviceID int, id int) (ConfigKey, error) {
	key, err := postgresutils.GetOne[ConfigKey](r.pool, ctx, getConfigKeyByIDSql, serviceID, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return key, ErrNotFound
	}

	return key, err
}

func (r *PostgresRepository) GetConfigKeyByName(ctx context.Context, serviceID int, name string) (ConfigKey, error) {
	key, err := postgresutils.GetOne[ConfigKey](r.pool, ctx, getConfigKeyByNameSql, serviceID, name)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return key, ErrNotFound
	}

	return key, err

}

func (r *PostgresRepository) ListConfigKeys(ctx context.Context, serviceIDs ...int) ([]ConfigKey, error) {
	if serviceIDs != nil {
		return postgresutils.GetAll[ConfigKey](r.pool, ctx, getAllConfigKeysByService, serviceIDs)
	} else {
		return postgresutils.GetAll[ConfigKey](r.pool, ctx, getAllConfigKeys)
	}
}
