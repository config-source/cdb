package postgres

import (
	"context"
	_ "embed"
	"errors"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/postgresutils"
	"github.com/jackc/pgx/v5"
)

//go:embed queries/configKeys/create_config_key.sql
var createConfigKeySql string

//go:embed queries/configKeys/get_config_key_by_id.sql
var getConfigKeyByIDSql string

//go:embed queries/configKeys/get_config_key_by_name.sql
var getConfigKeyByNameSql string

//go:embed queries/configKeys/get_all_config_keys.sql
var getAllConfigKeys string

func (r *Repository) CreateConfigKey(ctx context.Context, ck cdb.ConfigKey) (cdb.ConfigKey, error) {
	var canPropagate bool
	if ck.CanPropagate == nil {
		canPropagate = true
	} else {
		canPropagate = *ck.CanPropagate
	}

	return postgresutils.GetOne[cdb.ConfigKey](r.pool, ctx, createConfigKeySql, ck.Name, ck.ValueType, canPropagate)
}

func (r *Repository) GetConfigKey(ctx context.Context, id int) (cdb.ConfigKey, error) {
	key, err := postgresutils.GetOne[cdb.ConfigKey](r.pool, ctx, getConfigKeyByIDSql, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return key, cdb.ErrConfigKeyNotFound
	}

	return key, err
}

func (r *Repository) GetConfigKeyByName(ctx context.Context, name string) (cdb.ConfigKey, error) {
	key, err := postgresutils.GetOne[cdb.ConfigKey](r.pool, ctx, getConfigKeyByNameSql, name)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return key, cdb.ErrConfigKeyNotFound
	}

	return key, err

}

func (r *Repository) ListConfigKeys(ctx context.Context) ([]cdb.ConfigKey, error) {
	return postgresutils.GetAll[cdb.ConfigKey](r.pool, ctx, getAllConfigKeys)
}
