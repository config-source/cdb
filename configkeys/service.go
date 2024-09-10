package configkeys

import (
	"context"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/repository"
)

type Service struct {
	auth auth.AuthorizationGateway
	repo repository.ModelRepository
}

func NewService(repo repository.ModelRepository, auth auth.AuthorizationGateway) *Service {
	return &Service{
		auth: auth,
		repo: repo,
	}
}

func (svc *Service) CreateConfigKey(ctx context.Context, actor auth.User, env cdb.ConfigKey) (cdb.ConfigKey, error) {
	canManageConfigKeys, err := svc.auth.HasPermission(ctx, actor, auth.PermissionManageConfigKeys)
	if err != nil {
		return cdb.ConfigKey{}, err
	}

	if !canManageConfigKeys {
		return cdb.ConfigKey{}, auth.ErrUnauthorized
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

func (svc *Service) GetConfigKeyByName(ctx context.Context, actor auth.User, name string) (cdb.ConfigKey, error) {
	authErr := svc.hasReadPermissions(ctx, actor)
	if authErr != nil {
		return cdb.ConfigKey{}, authErr
	}

	return svc.repo.GetConfigKeyByName(ctx, name)
}

func (svc *Service) GetConfigKeyByID(ctx context.Context, actor auth.User, id int) (cdb.ConfigKey, error) {
	authErr := svc.hasReadPermissions(ctx, actor)
	if authErr != nil {
		return cdb.ConfigKey{}, authErr
	}

	return svc.repo.GetConfigKey(ctx, id)
}

func (svc *Service) ListConfigKeys(ctx context.Context, actor auth.User) ([]cdb.ConfigKey, error) {
	authErr := svc.hasReadPermissions(ctx, actor)
	if authErr != nil {
		return nil, authErr
	}

	return svc.repo.ListConfigKeys(ctx)
}