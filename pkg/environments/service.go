package environments

import (
	"context"

	"github.com/config-source/cdb/pkg/auth"
)

type Service struct {
	auth auth.AuthorizationGateway
	repo *Repository
}

func NewService(repo *Repository, auth auth.AuthorizationGateway) *Service {
	return &Service{
		auth: auth,
		repo: repo,
	}
}

func (svc *Service) CreateEnvironment(ctx context.Context, actor auth.User, env Environment) (Environment, error) {
	canManageEnvironments, err := svc.auth.HasPermission(ctx, actor, auth.PermissionManageEnvironments)
	if err != nil {
		return Environment{}, err
	}

	if !canManageEnvironments {
		return Environment{}, auth.ErrUnauthorized
	}

	return svc.repo.CreateEnvironment(ctx, env)
}

func (svc *Service) singleRetrievalPermissionChecks(ctx context.Context, actor auth.User, env Environment, retrievalErr error) (Environment, error) {
	canManageEnvironments, err := svc.auth.HasPermission(ctx, actor, auth.PermissionManageEnvironments)
	if err != nil {
		return Environment{}, err
	}

	canConfigureSensitiveEnvironments, err := svc.auth.HasPermission(
		ctx,
		actor,
		auth.PermissionConfigureSensitiveEnvironments,
	)
	if err != nil {
		return Environment{}, err
	}

	canConfigureEnvironments, err := svc.auth.HasPermission(
		ctx,
		actor,
		auth.PermissionConfigureEnvironments,
	)
	if err != nil {
		return Environment{}, err
	}

	if !canManageEnvironments && !canConfigureEnvironments && !canConfigureSensitiveEnvironments {
		return Environment{}, auth.ErrUnauthorized
	}

	if !(canConfigureSensitiveEnvironments || canManageEnvironments) && env.Sensitive {
		return Environment{}, ErrNotFound
	}

	return env, retrievalErr
}

func (svc *Service) GetEnvironmentByName(ctx context.Context, actor auth.User, serviceName, name string) (Environment, error) {
	env, err := svc.repo.GetEnvironmentByName(ctx, serviceName, name)
	return svc.singleRetrievalPermissionChecks(
		ctx,
		actor,
		env,
		err,
	)
}

func (svc *Service) GetEnvironmentByID(ctx context.Context, actor auth.User, id int) (Environment, error) {
	env, err := svc.repo.GetEnvironment(ctx, id)
	return svc.singleRetrievalPermissionChecks(
		ctx,
		actor,
		env,
		err,
	)
}

func (svc *Service) ListEnvironments(ctx context.Context, actor auth.User) ([]Environment, error) {
	canManageEnvironments, err := svc.auth.HasPermission(ctx, actor, auth.PermissionManageEnvironments)
	if err != nil {
		return nil, err
	}

	if canManageEnvironments {
		return svc.repo.ListEnvironments(ctx, true)
	}

	canSeeEnvirons, err := svc.auth.HasPermission(ctx, actor, auth.PermissionConfigureEnvironments)
	if err != nil {
		return nil, err
	}

	canSeeSensitiveEnvirons, err := svc.auth.HasPermission(ctx, actor, auth.PermissionConfigureSensitiveEnvironments)
	if err != nil {
		return nil, err
	}

	if !canSeeEnvirons && !canSeeSensitiveEnvirons {
		return nil, auth.ErrUnauthorized
	}

	return svc.repo.ListEnvironments(ctx, canSeeSensitiveEnvirons)
}

func getChildren(parent Environment, environments []Environment) []Tree {
	children := []Tree{}

	for _, env := range environments {
		isChild := env.PromotesToID != nil && *env.PromotesToID == parent.ID
		if isChild {
			children = append(children, Tree{
				Environment: env,
				Children:    getChildren(env, environments),
			})
		}
	}

	return children
}

func (svc *Service) GetEnvironmentTree(ctx context.Context, actor auth.User) ([]Tree, error) {
	canManageEnvironments, err := svc.auth.HasPermission(ctx, actor, auth.PermissionManageEnvironments)
	if err != nil {
		return nil, err
	}

	if !canManageEnvironments {
		return nil, auth.ErrUnauthorized
	}

	environs, err := svc.repo.ListEnvironments(ctx, true)
	if err != nil {
		return nil, err
	}

	trees := []Tree{}
	for _, env := range environs {
		if env.PromotesToID == nil {
			trees = append(trees, Tree{
				Environment: env,
				Children:    getChildren(env, environs),
			})
		}
	}

	return trees, nil
}

func (svc *Service) UpdateEnvironment(ctx context.Context, actor auth.User, env Environment) (Environment, error) {
	canManageEnvironments, err := svc.auth.HasPermission(ctx, actor, auth.PermissionManageEnvironments)
	if err != nil {
		return Environment{}, err
	}

	if !canManageEnvironments {
		return Environment{}, auth.ErrUnauthorized
	}

	updated, err := svc.repo.UpdateEnvironment(ctx, env)
	updated.Service = env.Service
	return updated, err
}

func (svc *Service) DeleteEnvironment(ctx context.Context, actor auth.User, id int) error {
	canManageEnvironments, err := svc.auth.HasPermission(ctx, actor, auth.PermissionManageEnvironments)
	if err != nil {
		return err
	}

	if !canManageEnvironments {
		return auth.ErrUnauthorized
	}

	return svc.repo.DeleteEnvironment(ctx, id)
}
