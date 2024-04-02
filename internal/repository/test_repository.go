package repository

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/config-source/cdb"
)

// TestRepository is an in-memory ModelRepository used for tests only.
type TestRepository struct {
	IsHealthy bool
	Error     error

	Environments map[int]cdb.Environment
	ConfigKeys   map[int]cdb.ConfigKey
	ConfigValues map[int]cdb.ConfigValue
}

func (tr *TestRepository) Healthy(ctx context.Context) bool {
	return tr.IsHealthy
}

func (tr *TestRepository) GetEnvironmentByName(ctx context.Context, name string) (cdb.Environment, error) {
	if tr.Error != nil {
		return cdb.Environment{}, tr.Error
	}

	for _, env := range tr.Environments {
		if env.Name == name {
			return env, nil
		}
	}

	return cdb.Environment{}, cdb.ErrEnvNotFound
}

func (tr *TestRepository) GetEnvironment(ctx context.Context, id int) (cdb.Environment, error) {
	if tr.Error != nil {
		return cdb.Environment{}, tr.Error
	}

	if env, ok := tr.Environments[id]; ok {
		return env, nil
	}

	return cdb.Environment{}, cdb.ErrEnvNotFound
}

func (tr *TestRepository) CreateEnvironment(ctx context.Context, env cdb.Environment) (cdb.Environment, error) {
	if tr.Error != nil {
		return cdb.Environment{}, tr.Error
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

func (tr *TestRepository) CreateConfigKey(ctx context.Context, ck cdb.ConfigKey) (cdb.ConfigKey, error) {
	ck.ID = len(tr.ConfigKeys) + 1
	ck.CreatedAt = time.Now()
	tr.ConfigKeys[ck.ID] = ck
	return ck, nil
}

func (tr *TestRepository) GetConfigKey(ctx context.Context, id int) (cdb.ConfigKey, error) {
	if tr.Error != nil {
		return cdb.ConfigKey{}, tr.Error
	}

	if ck, ok := tr.ConfigKeys[id]; ok {
		return ck, nil
	}

	return cdb.ConfigKey{}, cdb.ErrConfigKeyNotFound
}

func (tr *TestRepository) ListConfigKeys(ctx context.Context) ([]cdb.ConfigKey, error) {
	keys := make([]cdb.ConfigKey, len(tr.ConfigKeys))

	for id, ck := range tr.ConfigKeys {
		keys[id-1] = ck
	}

	return keys, nil
}

func (tr *TestRepository) CreateConfigValue(ctx context.Context, cv cdb.ConfigValue) (cdb.ConfigValue, error) {
	_, err := tr.GetEnvironment(ctx, cv.EnvironmentID)
	if err != nil {
		return cv, err
	}

	_, err = tr.GetConfigKey(ctx, cv.ConfigKeyID)
	if err != nil {
		return cv, err
	}

	cv.ID = len(tr.ConfigKeys) + 1
	cv.CreatedAt = time.Now()
	tr.ConfigValues[cv.ID] = cv
	return cv, nil
}

func keyAlreadyInSet(values []cdb.ConfigValue, newValue cdb.ConfigValue) bool {
	for _, cv := range values {
		if cv.ID == newValue.ID {
			return true
		}
	}

	return false
}

func (tr *TestRepository) GetConfiguration(ctx context.Context, environmentID int) ([]cdb.ConfigValue, error) {
	values := make([]cdb.ConfigValue, 0)

	env, err := tr.GetEnvironment(ctx, environmentID)
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
			values = append(values, cv)
		}
	}

	if env.PromotesToID != nil {
		parentValues, err := tr.GetConfiguration(ctx, *env.PromotesToID)
		if err != nil {
			return values, err
		}

		for _, cv := range parentValues {
			if !keyAlreadyInSet(values, cv) {
				values = append(values, cv)
			}
		}
	}

	return values, nil
}

func (tr *TestRepository) GetConfigurationValue(ctx context.Context, environmentName, key string) (cdb.ConfigValue, error) {
	env, err := tr.GetEnvironmentByName(ctx, environmentName)
	if err != nil {
		return cdb.ConfigValue{}, err
	}

	validEnvIds := []int{env.ID}
	for env.PromotesToID != nil {
		validEnvIds = append(validEnvIds, *env.PromotesToID)

		env, err = tr.GetEnvironment(ctx, *env.PromotesToID)
		if err != nil {
			return cdb.ConfigValue{}, err
		}
	}

	// Gotta be a smarter way to do this...
	var bestMatch cdb.ConfigValue
	var bestMatchIndex int = math.MaxInt

	for _, cv := range tr.ConfigValues {
		matchedIndex := -1

		for idx, id := range validEnvIds {
			if cv.EnvironmentID == id {
				matchedIndex = idx
				break
			}
		}

		if matchedIndex == -1 {
			continue
		}

		ck, err := tr.GetConfigKey(ctx, cv.ConfigKeyID)
		if err != nil {
			return cdb.ConfigValue{}, err
		}

		if ck.Name == key && bestMatchIndex > matchedIndex {
			bestMatch = cv
			bestMatchIndex = matchedIndex
		}
	}

	if bestMatchIndex == math.MaxInt {
		return cdb.ConfigValue{}, cdb.ErrConfigValueNotFound
	}

	return bestMatch, nil
}