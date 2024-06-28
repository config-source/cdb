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

func isUniqueConstraint(err error) bool {
	return strings.Contains(err.Error(), "unique constraint")
}

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
	if err != nil && isUniqueConstraint(err) {
		return nil, cdb.ErrConfigValueAlreadySet
	} else if err != nil {
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

func (r *Repository) GetConfigValueByEnvAndKey(ctx context.Context, environmentName string, key string) (*cdb.ConfigValue, error) {
	cv, err := getOne[cdb.ConfigValue](
		r,
		ctx,
		getConfigValueByEnvironmentAndKeySql,
		environmentName,
		key,
	)
	if errors.Is(err, pgx.ErrNoRows) {
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

func (r *Repository) getPromotesToName(ctx context.Context, envName string) (*string, error) {
	var promotesToName *string
	err := r.Raw().QueryRow(
		ctx,
		"SELECT name FROM environments WHERE id = (SELECT promotes_to_id FROM environments WHERE name = $1)",
		envName,
	).Scan(&promotesToName)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return promotesToName, err
	}

	return promotesToName, nil
}

func getConfigurationRecursively(ctx context.Context, r *Repository, environmentName string, excludedKeys []string) ([]cdb.ConfigValue, error) {
	immediateValues, err := getAll[cdb.ConfigValue](r, ctx, getAllConfigValuesForEnvironmentExceptKeysSql, environmentName, excludedKeys)
	if err != nil {
		return immediateValues, err
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

func (r *Repository) GetConfiguration(ctx context.Context, environmentName string) ([]cdb.ConfigValue, error) {
	immediateValues, err := getAll[cdb.ConfigValue](r, ctx, getAllConfigValuesForEnvironmentSql, environmentName)
	if err != nil {
		return immediateValues, err
	}

	promotesToName, _ := r.getPromotesToName(ctx, environmentName)
	if promotesToName != nil {
		parentValues, err := getConfigurationRecursively(ctx, r, *promotesToName, getAllKeys(immediateValues))
		for idx := range parentValues {
			parentValues[idx].Inherited = true
			parentValues[idx].InheritedFrom = *promotesToName
		}

		return append(immediateValues, parentValues...), err
	}

	return immediateValues, nil
}

func (r *Repository) GetConfigurationValue(ctx context.Context, environmentName, key string) (*cdb.ConfigValue, error) {
	cv, err := r.GetConfigValueByEnvAndKey(ctx, environmentName, key)
	if errors.Is(err, cdb.ErrConfigValueNotFound) {
		promotesToName, _ := r.getPromotesToName(ctx, environmentName)
		if promotesToName != nil {
			cv, err := r.GetConfigurationValue(ctx, *promotesToName, key)
			cv.Inherited = true
			cv.InheritedFrom = *promotesToName
			return cv, err
		}

		return cv, cdb.ErrConfigValueNotFound
	}

	return cv, err
}

func (r *Repository) GetConfigurationValueByID(ctx context.Context, configValueID int) (*cdb.ConfigValue, error) {
	cv, err := getOne[cdb.ConfigValue](r, ctx, getConfigValueByIDSql, configValueID)
	return &cv, err
}
