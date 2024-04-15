package postgres

import (
	"context"
	_ "embed"
	"errors"
	"strings"

	"github.com/config-source/cdb"
	"github.com/jackc/pgx/v5"
)

//go:embed queries/configValues/create_config_value.sql
var createConfigValueSql string

//go:embed queries/configValues/update_config_value.sql
var updateConfigValueSql string

//go:embed queries/configValues/get_config_value_by_id.sql
var getConfigValueByIDSql string

//go:embed queries/configValues/get_config_value_by_environment_and_key.sql
var getConfigValueByEnvironmentAndKeySql string

//go:embed queries/configValues/get_all_config_values_for_environment.sql
var getAllConfigValuesForEnvironmentSql string

//go:embed queries/configValues/get_all_config_values_except_matching_keys.sql
var getAllConfigValuesForEnvironmentExceptKeysSql string

func (r *Repository) CreateConfigValue(ctx context.Context, cv *cdb.ConfigValue) (*cdb.ConfigValue, error) {
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
	if err != nil {
		return nil, err
	}

	created, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[cdb.ConfigValue])
	return &created, err
}

func (r *Repository) UpdateConfigurationValue(ctx context.Context, cv *cdb.ConfigValue) (*cdb.ConfigValue, error) {
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

	updated, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[cdb.ConfigValue])
	if err == nil {
		return &updated, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return &updated, cdb.ErrConfigValueNotFound
	}

	msg := err.Error()
	if strings.Contains(msg, "config_values_environment_id_fkey") {
		return &updated, cdb.ErrEnvNotFound
	}

	if strings.Contains(msg, "config_values_config_key_id_fkey") {
		return &updated, cdb.ErrConfigKeyNotFound
	}

	return &updated, err
}

func (r *Repository) GetConfigValue(ctx context.Context, environmentID int, key string) (*cdb.ConfigValue, error) {
	cv, err := getOne[cdb.ConfigValue](
		r,
		ctx,
		getConfigValueByEnvironmentAndKeySql,
		environmentID,
		key,
	)
	if err == nil {
		return &cv, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		promotesToID, _ := r.getPromotesToID(ctx, environmentID)
		if promotesToID != nil {
			cv, err := r.GetConfigValue(ctx, *promotesToID, key)
			cv.Inherited = true
			return cv, err
		}

		return &cv, cdb.ErrConfigValueNotFound
	}

	return &cv, err
}

func getAllKeys(values []cdb.ConfigValue) []string {
	keys := make([]string, len(values))
	for idx, cv := range values {
		keys[idx] = cv.Name
	}
	return keys
}

func (r *Repository) getPromotesToID(ctx context.Context, envID int) (*int, error) {
	var promotesToID *int
	err := r.Raw().QueryRow(
		ctx,
		"SELECT promotes_to_id FROM environments WHERE id = $1",
		envID,
	).Scan(&promotesToID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return promotesToID, err
	}

	return promotesToID, nil
}

func getConfigurationRecursively(ctx context.Context, r *Repository, environmentID int, excludedKeys []string) ([]cdb.ConfigValue, error) {
	immediateValues, err := getAll[cdb.ConfigValue](r, ctx, getAllConfigValuesForEnvironmentExceptKeysSql, environmentID, excludedKeys)
	if err != nil {
		return immediateValues, err
	}

	promotesToID, err := r.getPromotesToID(ctx, environmentID)
	if err != nil {
		return immediateValues, err
	}

	if promotesToID != nil {
		parentValues, err := getConfigurationRecursively(ctx, r, *promotesToID, append(excludedKeys, getAllKeys(immediateValues)...))
		return append(immediateValues, parentValues...), err
	}

	return immediateValues, nil
}

func (r *Repository) GetConfiguration(ctx context.Context, environmentName string) ([]cdb.ConfigValue, error) {
	immediateValues, err := getAll[cdb.ConfigValue](r, ctx, getAllConfigValuesForEnvironmentSql, environmentName)
	if err != nil {
		return immediateValues, err
	}

	var promotesToID *int
	err = r.Raw().QueryRow(
		ctx,
		"SELECT promotes_to_id FROM environments WHERE name = $1",
		environmentName,
	).Scan(&promotesToID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return immediateValues, err
	}

	if promotesToID != nil {
		parentValues, err := getConfigurationRecursively(ctx, r, *promotesToID, getAllKeys(immediateValues))
		for idx := range parentValues {
			parentValues[idx].Inherited = true
		}

		return append(immediateValues, parentValues...), err
	}

	return immediateValues, nil
}

func (r *Repository) GetConfigurationValue(ctx context.Context, environmentName, key string) (*cdb.ConfigValue, error) {
	env, err := r.GetEnvironmentByName(ctx, environmentName)
	if err != nil {
		return nil, err
	}

	return r.GetConfigValue(ctx, env.ID, key)
}

func (r *Repository) GetConfigurationValueByID(ctx context.Context, configValueID int) (*cdb.ConfigValue, error) {
	cv, err := getOne[cdb.ConfigValue](r, ctx, getConfigValueByIDSql, configValueID)
	return &cv, err
}
