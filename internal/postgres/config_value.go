package postgres

import (
	"context"
	_ "embed"
	"errors"

	"github.com/config-source/cdb"
	"github.com/jackc/pgx/v5"
)

//go:embed queries/configValues/create_config_value.sql
var createConfigValueSql string

//go:embed queries/configValues/get_config_value_by_id.sql
var getConfigValueByIDSql string

//go:embed queries/configValues/get_config_value_by_environment_and_key.sql
var getConfigValueByEnvironmentAndKeySql string

//go:embed queries/configValues/get_all_config_values_for_environment.sql
var getAllConfigValuesForEnvironmentSql string

//go:embed queries/configValues/get_all_config_values_except_matching_keys.sql
var getAllConfigValuesForEnvironmentExceptKeysSql string

func (r *Repository) CreateConfigValue(ctx context.Context, cv cdb.ConfigValue) (cdb.ConfigValue, error) {
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
		var def cdb.ConfigValue
		return def, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[cdb.ConfigValue])
}

func (r *Repository) GetConfigValue(ctx context.Context, environmentID int, key string) (cdb.ConfigValue, error) {
	return getOne[cdb.ConfigValue](
		r,
		ctx,
		getConfigValueByEnvironmentAndKeySql,
		environmentID,
		key,
	)
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

func (r *Repository) GetConfigurationValue(ctx context.Context, environmentName, key string) (cdb.ConfigValue, error) {
	return getOne[cdb.ConfigValue](r, ctx, getConfigValueByEnvironmentAndKeySql, environmentName, key)
}

func (r *Repository) GetConfigurationValueByID(ctx context.Context, configValueID int) (cdb.ConfigValue, error) {
	return getOne[cdb.ConfigValue](r, ctx, getConfigValueByIDSql, configValueID)
}
