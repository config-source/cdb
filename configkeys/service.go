package configkeys

import (
	"context"

	"github.com/config-source/cdb/auth"
)

type Service struct {
	auth auth.AuthorizationGateway
	repo Repository
}

func NewService(repo Repository, auth auth.AuthorizationGateway) *Service {
	return &Service{
		auth: auth,
		repo: repo,
	}
}

func (svc *Service) CreateConfigKey(ctx context.Context, actor auth.User, env ConfigKey) (ConfigKey, error) {
	canManageConfigKeys, err := svc.auth.HasPermission(ctx, actor, auth.PermissionManageConfigKeys)
	if err != nil {
		return ConfigKey{}, err
	}

	if !canManageConfigKeys {
		return ConfigKey{}, auth.ErrUnauthorized
	}

	return svc.repo.CreateConfigKey(ctx, env)
}

func (svc *Service) hasReadPermissions(ctx context.Context, actor auth.User) error {
	canManageConfigKeys, err := svc.auth.HasPermission(ctx, actor, auth.PermissionManageConfigKeys)
	if err != nil {
		return err
	}

	canConfigureSensitiveEnviroments, err := svc.auth.HasPermission(
		ctx,
		actor,
		auth.PermissionConfigureSensitiveEnvironments,
	)
	if err != nil {
		return err
	}

	canConfigureEnvironments, err := svc.auth.HasPermission(
		ctx,
		actor,
		auth.PermissionConfigureEnvironments,
	)
	if err != nil {
		return err
	}

	if !canConfigureEnvironments && !canConfigureSensitiveEnviroments && !canManageConfigKeys {
		return auth.ErrUnauthorized
	}

	return nil
}

func (svc *Service) GetConfigKeyByName(ctx context.Context, actor auth.User, name string) (ConfigKey, error) {
	authErr := svc.hasReadPermissions(ctx, actor)
	if authErr != nil {
		return ConfigKey{}, authErr
	}

	return svc.repo.GetConfigKeyByName(ctx, name)
}

func (svc *Service) GetConfigKeyByID(ctx context.Context, actor auth.User, id int) (ConfigKey, error) {
	authErr := svc.hasReadPermissions(ctx, actor)
	if authErr != nil {
		return ConfigKey{}, authErr
	}

	return svc.repo.GetConfigKey(ctx, id)
}

func (svc *Service) ListConfigKeys(ctx context.Context, actor auth.User) ([]ConfigKey, error) {
	authErr := svc.hasReadPermissions(ctx, actor)
	if authErr != nil {
		return nil, authErr
	}

	return svc.repo.ListConfigKeys(ctx)
}
