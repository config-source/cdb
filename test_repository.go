package cdb

import (
	"context"
	"errors"
	"time"

	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/configvalues"
	"github.com/config-source/cdb/environments"
)

// TestRepository is an in-memory ModelRepository used for tests only.
type TestRepository struct {
	IsHealthy bool
	Error     error

	Environments map[int]environments.Environment
	ConfigKeys   map[int]configkeys.ConfigKey
	ConfigValues map[int]*configvalues.ConfigValue
}

func (tr *TestRepository) Healthy(ctx context.Context) bool {
	return tr.IsHealthy
}

func (tr *TestRepository) ListEnvironments(ctx context.Context, includeSensitive bool) ([]environments.Environment, error) {
	if tr.Error != nil {
		return nil, tr.Error
	}

	envs := make([]environments.Environment, len(tr.Environments))
	for id, env := range tr.Environments {
		if includeSensitive && env.Sensitive {
			envs[id-1] = env
		} else if env.Sensitive {
			continue
		} else {
			envs[id-1] = env
		}
	}

	return envs, nil
}

func (tr *TestRepository) GetEnvironmentByName(ctx context.Context, name string) (environments.Environment, error) {
	if tr.Error != nil {
		return environments.Environment{}, tr.Error
	}

	for _, env := range tr.Environments {
		if env.Name == name {
			return env, nil
		}
	}

	return environments.Environment{}, environments.ErrNotFound
}

func (tr *TestRepository) GetEnvironment(ctx context.Context, id int) (environments.Environment, error) {
	if tr.Error != nil {
		return environments.Environment{}, tr.Error
	}

	if env, ok := tr.Environments[id]; ok {
		return env, nil
	}

	return environments.Environment{}, environments.ErrNotFound
}

func (tr *TestRepository) CreateEnvironment(ctx context.Context, env environments.Environment) (environments.Environment, error) {
	if tr.Error != nil {
		return environments.Environment{}, tr.Error
	}

	if tr.Environments == nil {
		tr.Environments = make(map[int]environments.Environment)
	}

	env.ID = len(tr.Environments) + 1
	env.CreatedAt = time.Now()
	if env.PromotesToID != nil && *env.PromotesToID != 0 {
		if _, ok := tr.Environments[*env.PromotesToID]; !ok {
			return env, errors.New("promotes to id does not exist!")
		}
	}

	tr.Environments[env.ID] = env
	return env, nil
}

func (tr *TestRepository) CreateConfigKey(ctx context.Context, ck configkeys.ConfigKey) (configkeys.ConfigKey, error) {
	if tr.ConfigKeys == nil {
		tr.ConfigKeys = make(map[int]configkeys.ConfigKey)
	}

	ck.ID = len(tr.ConfigKeys) + 1
	ck.CreatedAt = time.Now()
	tr.ConfigKeys[ck.ID] = ck
	return ck, nil
}

func (tr *TestRepository) GetConfigKey(ctx context.Context, id int) (configkeys.ConfigKey, error) {
	if tr.Error != nil {
		return configkeys.ConfigKey{}, tr.Error
	}

	if ck, ok := tr.ConfigKeys[id]; ok {
		return ck, nil
	}

	return configkeys.ConfigKey{}, configkeys.ErrNotFound
}

func (tr *TestRepository) GetConfigKeyByName(ctx context.Context, name string) (configkeys.ConfigKey, error) {
	if tr.Error != nil {
		return configkeys.ConfigKey{}, tr.Error
	}

	for _, ck := range tr.ConfigKeys {
		if ck.Name == name {
			return ck, nil
		}
	}

	return configkeys.ConfigKey{}, configkeys.ErrNotFound
}

func (tr *TestRepository) ListConfigKeys(ctx context.Context) ([]configkeys.ConfigKey, error) {
	keys := make([]configkeys.ConfigKey, len(tr.ConfigKeys))

	for id, ck := range tr.ConfigKeys {
		keys[id-1] = ck
	}

	return keys, nil
}

func (tr *TestRepository) CreateConfigValue(ctx context.Context, cv *configvalues.ConfigValue) (*configvalues.ConfigValue, error) {
	_, err := tr.GetEnvironment(ctx, cv.EnvironmentID)
	if err != nil {
		return cv, err
	}

	_, err = tr.GetConfigKey(ctx, cv.ConfigKeyID)
	if err != nil {
		return cv, err
	}

	if tr.ConfigValues == nil {
		tr.ConfigValues = make(map[int]*configvalues.ConfigValue)
	}

	cv.ID = len(tr.ConfigValues) + 1
	cv.CreatedAt = time.Now()
	tr.ConfigValues[cv.ID] = cv
	return cv, nil
}

func (tr *TestRepository) UpdateConfigurationValue(ctx context.Context, cv *configvalues.ConfigValue) (*configvalues.ConfigValue, error) {
	_, err := tr.GetEnvironment(ctx, cv.EnvironmentID)
	if err != nil {
		return cv, err
	}

	_, err = tr.GetConfigKey(ctx, cv.ConfigKeyID)
	if err != nil {
		return cv, err
	}

	if tr.ConfigValues == nil {
		tr.ConfigValues = make(map[int]*configvalues.ConfigValue)
	}

	cv.CreatedAt = time.Now()
	tr.ConfigValues[cv.ID] = cv
	return cv, nil
}

func keyAlreadyInSet(values []configvalues.ConfigValue, newValue configvalues.ConfigValue) bool {
	for _, cv := range values {
		if cv.ConfigKeyID == newValue.ConfigKeyID {
			return true
		}
	}

	return false
}

func (tr *TestRepository) GetConfiguration(ctx context.Context, environmentName string) ([]configvalues.ConfigValue, error) {
	values := make([]configvalues.ConfigValue, 0)

	env, err := tr.GetEnvironmentByName(ctx, environmentName)
	if err != nil {
		return values, err
	}

	for _, cv := range tr.ConfigValues {
		if cv.EnvironmentID == env.ID {
			ck, err := tr.GetConfigKey(ctx, cv.ConfigKeyID)
			if err != nil {
				return values, err
			}

			cv.Name = ck.Name
			cv.ValueType = ck.ValueType
			values = append(values, *cv)
		}
	}

	if env.PromotesToID != nil {
		parent, err := tr.GetEnvironment(ctx, *env.PromotesToID)
		if err != nil {
			return values, err
		}

		parentValues, err := tr.GetConfiguration(ctx, parent.Name)
		if err != nil {
			return values, err
		}

		for _, cv := range parentValues {
			if !keyAlreadyInSet(values, cv) {
				cv.Inherited = true
				cv.InheritedFrom = parent.Name
				values = append(values, cv)
			}
		}
	}

	return values, nil
}

func (tr *TestRepository) GetConfigValueByEnvAndKey(ctx context.Context, environmentName, key string) (*configvalues.ConfigValue, error) {
	cv, err := tr.GetConfigurationValue(ctx, environmentName, key)
	if cv != nil && cv.Inherited {
		return nil, configvalues.ErrNotFound
	}

	return cv, err
}

func (tr *TestRepository) GetConfigurationValue(ctx context.Context, environmentName, key string) (*configvalues.ConfigValue, error) {
	configValues, err := tr.GetConfiguration(ctx, environmentName)
	if err != nil {
		return nil, err
	}

	for _, cv := range configValues {
		if cv.Name == key {
			return &cv, nil
		}
	}

	return nil, configvalues.ErrNotFound
}