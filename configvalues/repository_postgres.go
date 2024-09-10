package configvalues

import (
	"context"
	_ "embed"
	"errors"
	"strings"

	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/environments"
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

//go:embed queries/create_config_value.sql
var createConfigValueSql string

//go:embed queries/update_config_value.sql
var updateConfigValueSql string

//go:embed queries/get_config_value_by_id.sql
var getConfigValueByIDSql string

//go:embed queries/get_config_value_by_environment_and_key.sql
var getConfigValueByEnvironmentAndKeySql string

//go:embed queries/get_all_config_values_for_environment.sql
var getAllConfigValuesForEnvironmentSql string

//go:embed queries/get_all_config_values_except_matching_keys.sql
var getAllConfigValuesForEnvironmentExceptKeysSql string

func (r *Repository) CreateConfigValue(ctx context.Context, cv *ConfigValue) (*ConfigValue, error) {
	rows, err := r.pool.Query(
		ctx,
		createConfigValueSql,
		cv.EnvironmentID,
		cv.ConfigKeyID,
		cv.StrValue,
		cv.IntValue,
		cv.FloatValue,
		cv.BoolValue,
	)
	if err != nil && postgresutils.IsUniqueConstraintErr(err) {
		return nil, ErrAlreadySet
	} else if err != nil {
		return nil, err
	}

	created, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[ConfigValue])
	return &created, err
}

func (r *Repository) UpdateConfigurationValue(ctx context.Context, cv *ConfigValue) (*ConfigValue, error) {
	rows, err := r.pool.Query(
		ctx,
		updateConfigValueSql,
		cv.EnvironmentID,
		cv.ConfigKeyID,
		cv.StrValue,
		cv.IntValue,
		cv.FloatValue,
		cv.BoolValue,
		cv.ID,
	)
	if err != nil {
		return nil, err
	}

	updated, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[ConfigValue])
	if err == nil {
		return &updated, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return &updated, ErrNotFound
	}

	msg := err.Error()
	if strings.Contains(msg, "config_values_environment_id_fkey") {
		return &updated, environments.ErrNotFound
	}

	if strings.Contains(msg, "config_values_config_key_id_fkey") {
		return &updated, configkeys.ErrNotFound
	}

	return &updated, err
}

func (r *Repository) GetConfigValueByEnvAndKey(ctx context.Context, environmentName string, key string) (*ConfigValue, error) {
	cv, err := postgresutils.GetOne[ConfigValue](
		r.pool,
		ctx,
		getConfigValueByEnvironmentAndKeySql,
		environmentName,
		key,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return &cv, ErrNotFound
	}

	return &cv, err
}

func getAllKeys(values []ConfigValue) []string {
	keys := make([]string, len(values))
	for idx, cv := range values {
		keys[idx] = cv.Name
	}
	return keys
}

func (r *Repository) getPromotesToName(ctx context.Context, envName string) (*string, error) {
	var promotesToName *string
	err := r.pool.QueryRow(
		ctx,
		"SELECT name FROM environments WHERE id = (SELECT promotes_to_id FROM environments WHERE name = $1)",
		envName,
	).Scan(&promotesToName)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return promotesToName, err
	}

	return promotesToName, nil
}

func getConfigurationRecursively(ctx context.Context, r *Repository, environmentName string, excludedKeys []string) ([]ConfigValue, error) {
	immediateValues, err := postgresutils.GetAll[ConfigValue](r.pool, ctx, getAllConfigValuesForEnvironmentExceptKeysSql, environmentName, excludedKeys)
	if err != nil {
		return immediateValues, err
	}

	for idx := range immediateValues {
		immediateValues[idx].Inherited = true
		immediateValues[idx].InheritedFrom = environmentName
	}

	promotesToName, err := r.getPromotesToName(ctx, environmentName)
	if err != nil {
		return immediateValues, err
	}

	if promotesToName != nil {
		parentValues, err := getConfigurationRecursively(ctx, r, *promotesToName, append(excludedKeys, getAllKeys(immediateValues)...))
		return append(immediateValues, parentValues...), err
	}

	return immediateValues, nil
}

func (r *Repository) GetConfiguration(ctx context.Context, environmentName string) ([]ConfigValue, error) {
	immediateValues, err := postgresutils.GetAll[ConfigValue](r.pool, ctx, getAllConfigValuesForEnvironmentSql, environmentName)
	if err != nil {
		return immediateValues, err
	}

	promotesToName, _ := r.getPromotesToName(ctx, environmentName)
	if promotesToName != nil {
		parentValues, err := getConfigurationRecursively(ctx, r, *promotesToName, getAllKeys(immediateValues))
		return append(immediateValues, parentValues...), err
	}

	return immediateValues, nil
}

func (r *Repository) GetConfigurationValue(ctx context.Context, environmentName, key string) (*ConfigValue, error) {
	cv, err := r.GetConfigValueByEnvAndKey(ctx, environmentName, key)
	if errors.Is(err, ErrNotFound) {
		promotesToName, _ := r.getPromotesToName(ctx, environmentName)
		if promotesToName != nil {
			cv, err := r.GetConfigurationValue(ctx, *promotesToName, key)
			cv.Inherited = true
			cv.InheritedFrom = *promotesToName
			return cv, err
		}

		return cv, ErrNotFound
	}

	return cv, err
}

func (r *Repository) GetConfigurationValueByID(ctx context.Context, configValueID int) (*ConfigValue, error) {
	cv, err := postgresutils.GetOne[ConfigValue](r.pool, ctx, getConfigValueByIDSql, configValueID)
	return &cv, err
}
