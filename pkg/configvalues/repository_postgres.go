package configvalues

import (
	"context"
	_ "embed"
	"errors"
	"strings"

	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/postgresutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PostgresRepository struct {
	pool    *pgxpool.Pool
	log     zerolog.Logger
	envRepo environments.Repository
}

func NewRepository(log zerolog.Logger, pool *pgxpool.Pool, envRepo environments.Repository) Repository {
	return &PostgresRepository{
		log:     log,
		pool:    pool,
		envRepo: envRepo,
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

func (r *PostgresRepository) CreateConfigValue(ctx context.Context, cv *ConfigValue) (*ConfigValue, error) {
	created, err := postgresutils.GetOneLax[ConfigValue](
		r.pool,
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

	return &created, err
}

func (r *PostgresRepository) UpdateConfigurationValue(ctx context.Context, cv *ConfigValue) (*ConfigValue, error) {
	updated, err := postgresutils.GetOneLax[ConfigValue](
		r.pool,
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
	}

	return &updated, err
}

func (r *PostgresRepository) GetConfigValueByEnvAndKey(ctx context.Context, environmentID int, key string) (*ConfigValue, error) {
	cv, err := postgresutils.GetOne[ConfigValue](
		r.pool,
		ctx,
		getConfigValueByEnvironmentAndKeySql,
		environmentID,
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

func getConfigurationRecursively(ctx context.Context, r *PostgresRepository, environmentID int, excludedKeys []string) ([]ConfigValue, error) {
	env, err := r.envRepo.GetEnvironment(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	immediateValues, err := postgresutils.GetAll[ConfigValue](r.pool, ctx, getAllConfigValuesForEnvironmentExceptKeysSql, environmentID, excludedKeys)
	if err != nil {
		return immediateValues, err
	}

	for idx := range immediateValues {
		immediateValues[idx].Inherited = true
		immediateValues[idx].InheritedFrom = env.Name
	}

	if env.PromotesToID != nil {
		parentValues, err := getConfigurationRecursively(ctx, r, *env.PromotesToID, append(excludedKeys, getAllKeys(immediateValues)...))
		return append(immediateValues, parentValues...), err
	}

	return immediateValues, nil
}

func (r *PostgresRepository) GetConfiguration(ctx context.Context, environmentID int) ([]ConfigValue, error) {
	immediateValues, err := postgresutils.GetAll[ConfigValue](r.pool, ctx, getAllConfigValuesForEnvironmentSql, environmentID)
	if err != nil {
		return immediateValues, err
	}

	env, err := r.envRepo.GetEnvironment(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	if env.PromotesToID != nil {
		parentValues, err := getConfigurationRecursively(ctx, r, *env.PromotesToID, getAllKeys(immediateValues))
		return append(immediateValues, parentValues...), err
	}

	return immediateValues, nil
}

func (r *PostgresRepository) GetConfigurationValue(ctx context.Context, environmentID int, key string) (*ConfigValue, error) {
	cv, err := r.GetConfigValueByEnvAndKey(ctx, environmentID, key)
	if errors.Is(err, ErrNotFound) {
		env, err := r.envRepo.GetEnvironment(ctx, environmentID)
		if err != nil {
			return nil, err
		}

		if env.PromotesToID != nil {
			parent, err := r.envRepo.GetEnvironment(ctx, *env.PromotesToID)
			cv, err := r.GetConfigurationValue(ctx, parent.ID, key)
			cv.Inherited = true
			cv.InheritedFrom = parent.Name
			return cv, err
		}

		return cv, ErrNotFound
	}

	return cv, err
}

func (r *PostgresRepository) GetConfigurationValueByID(ctx context.Context, configValueID int) (*ConfigValue, error) {
	cv, err := postgresutils.GetOne[ConfigValue](r.pool, ctx, getConfigValueByIDSql, configValueID)
	return &cv, err
}
