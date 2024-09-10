package configvalues

import (
	"context"
	"errors"
	"fmt"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/repository"
)

var (
	ErrValueTypeMustBeSet = errors.New("must set ValueType on the config value when trying to dynamically create a key")
)

type Service struct {
	DynamicConfigKeys bool

	repo repository.ModelRepository
	auth auth.AuthorizationGateway
}

func NewService(
	repo repository.ModelRepository,
	auth auth.AuthorizationGateway,
	dynamicConfigKeys bool,
) *Service {
	return &Service{
		repo:              repo,
		auth:              auth,
		DynamicConfigKeys: dynamicConfigKeys,
	}
}

func (svc *Service) canConfigureEnvironment(
	ctx context.Context,
	actor auth.User,
	env cdb.Environment,
) error {
	canConfigure, err := svc.auth.HasPermission(ctx, actor, auth.PermissionConfigureEnvironments)
	if err != nil {
		return err
	}

	canConfigureSensitive, err := svc.auth.HasPermission(ctx, actor, auth.PermissionConfigureEnvironments)
	if err != nil {
		return err
	}

	if !canConfigureSensitive && env.Sensitive {
		return auth.ErrUnauthorized
	}

	if canConfigure || canConfigureSensitive {
		return nil
	}

	return auth.ErrUnauthorized
}

func (svc *Service) SetConfigurationValue(
	ctx context.Context,
	actor auth.User,
	envName string,
	key string,
	cv *cdb.ConfigValue,
) (*cdb.ConfigValue, error) {
	env, err := svc.repo.GetEnvironmentByName(ctx, envName)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment by name: %s", err)
	}

	if authErr := svc.canConfigureEnvironment(ctx, actor, env); authErr != nil {
		return nil, authErr
	}

	cv.EnvironmentID = env.ID

	ck, err := svc.repo.GetConfigKeyByName(ctx, key)
	shouldCreate := errors.Is(err, cdb.ErrConfigKeyNotFound) && svc.DynamicConfigKeys
	if shouldCreate {
		if cv.ValueType == 0 {
			return cv, ErrValueTypeMustBeSet
		}

		ck = cdb.NewConfigKey(key, cv.ValueType)
		ck, err = svc.repo.CreateConfigKey(ctx, ck)
		if err != nil {
			return cv, fmt.Errorf("failed to create new config key: %w", err)
		}

	} else if err != nil {
		return nil, fmt.Errorf("unable to retrieve config key: %w", err)
	}

	cv.ConfigKeyID = ck.ID

	// Force the ValueType to match the config key so that the client can't send
	// us a new ValueType that doesn't match it's config key thereby bypassing
	// the validity check.
	cv.ValueType = ck.ValueType
	if err := cv.Valid(); err != nil {
		return nil, err
	}

	var result *cdb.ConfigValue
	alreadySet, err := svc.repo.GetConfigValueByEnvAndKey(ctx, envName, key)
	if err != nil {
		result, err = svc.repo.CreateConfigValue(ctx, cv)
	} else {
		cv.ID = alreadySet.ID
		result, err = svc.repo.UpdateConfigurationValue(ctx, cv)
	}

	// Create and Update ConfigValue do not always populate these.
	result.ValueType = ck.ValueType
	result.Name = ck.Name

	return result, err
}

func (svc *Service) CreateConfigValue(
	ctx context.Context,
	actor auth.User,
	cv cdb.ConfigValue,
) (cdb.ConfigValue, error) {
	env, err := svc.repo.GetEnvironment(ctx, cv.EnvironmentID)
	if err != nil {
		return cdb.ConfigValue{}, err
	}

	if authErr := svc.canConfigureEnvironment(ctx, actor, env); authErr != nil {
		return cdb.ConfigValue{}, authErr
	}

	ck, err := svc.repo.GetConfigKey(ctx, cv.ConfigKeyID)
	if err != nil {
		return cdb.ConfigValue{}, err
	}

	cv.ValueType = ck.ValueType
	if err := cv.Valid(); err != nil {
		return cdb.ConfigValue{}, err
	}

	created, err := svc.repo.CreateConfigValue(ctx, &cv)
	if err != nil {
		return cdb.ConfigValue{}, err
	}

	return *created, nil
}

// TODO: there isn't really a concept of permissions for reading configuration
// yet and it's unclear if there ever will be. But these methods exist to future
// proof for it and to simplify things so that the API struct never needs to
// talk to ModelRepository directly.

func (svc *Service) GetConfiguration(ctx context.Context, actor auth.User, envName string) ([]cdb.ConfigValue, error) {
	return svc.repo.GetConfiguration(ctx, envName)
}

func (svc *Service) GetConfigurationValue(ctx context.Context, actor auth.User, envName, key string) (*cdb.ConfigValue, error) {
	return svc.repo.GetConfigurationValue(ctx, envName, key)
}
