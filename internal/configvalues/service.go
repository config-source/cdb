package configvalues

import (
	"context"
	"errors"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/internal/repository"
)

var (
	ErrValueTypeMustBeSet = errors.New("must set ValueType on the config value when trying to dynamically create a key")
)

type Service struct {
	DynamicConfigKeys bool

	repo repository.ModelRepository
}

func NewService(
	repo repository.ModelRepository,
	dynamicConfigKeys bool,
) *Service {
	return &Service{
		repo:              repo,
		DynamicConfigKeys: dynamicConfigKeys,
	}
}

func (s *Service) SetConfigurationValue(
	ctx context.Context,
	envName string,
	key string,
	cv cdb.ConfigValue,
) (cdb.ConfigValue, error) {
	env, err := s.repo.GetEnvironmentByName(ctx, envName)
	if err != nil {
		return cdb.ConfigValue{}, err
	}
	cv.EnvironmentID = env.ID

	ck, err := s.repo.GetConfigKeyByName(ctx, key)
	shouldCreate := err == cdb.ErrConfigKeyNotFound && s.DynamicConfigKeys
	if shouldCreate {
		if cv.ValueType == 0 {
			return cv, ErrValueTypeMustBeSet
		}

		ck = cdb.NewConfigKey(key, cv.ValueType)
		ck, err = s.repo.CreateConfigKey(ctx, ck)
		if err != nil {
			return cv, err
		}
	} else if err != nil {
		return cdb.ConfigValue{}, err
	}
	cv.ConfigKeyID = ck.ID
	// Force the ValueType to match the config key so that the client can't send
	// us a new ValueType that doesn't match it's config key thereby bypassing
	// the validity check.
	cv.ValueType = ck.ValueType
	if err := cv.Valid(); err != nil {
		return cdb.ConfigValue{}, err
	}

	created, err := s.repo.CreateConfigValue(ctx, &cv)
	if err != nil {
		return cdb.ConfigValue{}, err
	}

	return *created, nil
}

func (s *Service) CreateConfigValue(ctx context.Context, cv cdb.ConfigValue) (cdb.ConfigValue, error) {
	ck, err := s.repo.GetConfigKey(ctx, cv.ConfigKeyID)
	if err != nil {
		return cdb.ConfigValue{}, err
	}

	cv.ValueType = ck.ValueType
	if err := cv.Valid(); err != nil {
		return cdb.ConfigValue{}, err
	}

	created, err := s.repo.CreateConfigValue(ctx, &cv)
	if err != nil {
		return cdb.ConfigValue{}, err
	}

	return *created, nil
}
