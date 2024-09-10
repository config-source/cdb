package services

import (
	"context"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/repository"
)

type Environments struct {
	auth auth.AuthorizationGateway
	repo repository.ModelRepository
}

func NewEnvironmentsService(repo repository.ModelRepository, auth auth.AuthorizationGateway) *Environments {
	return &Environments{
		auth: auth,
		repo: repo,
	}
}

func (svc *Environments) CreateEnvironment(ctx context.Context, actor auth.User, env cdb.Environment) (cdb.Environment, error) {
	canManageEnvironments, err := svc.auth.HasPermission(ctx, actor, auth.PermissionManageEnvironments)
	if err != nil {
		return cdb.Environment{}, err
	}

	if !canManageEnvironments {
		return cdb.Environment{}, auth.ErrUnauthorized
	}

	return svc.repo.CreateEnvironment(ctx, env)
}

func (svc *Environments) singleRetrievalPermissionChecks(ctx context.Context, actor auth.User, env cdb.Environment, retrievalErr error) (cdb.Environment, error) {
	canManageEnvironments, err := svc.auth.HasPermission(ctx, actor, auth.PermissionManageEnvironments)
	if err != nil {
		return cdb.Environment{}, err
	}

	canConfigureSensitiveEnvironments, err := svc.auth.HasPermission(
		ctx,
		actor,
		auth.PermissionConfigureSensitiveEnvironments,
	)
	if err != nil {
		return cdb.Environment{}, err
	}

	canConfigureEnvironments, err := svc.auth.HasPermission(
		ctx,
		actor,
		auth.PermissionConfigureEnvironments,
	)
	if err != nil {
		return cdb.Environment{}, err
	}

	if !canManageEnvironments && !canConfigureEnvironments && !canConfigureSensitiveEnvironments {
		return cdb.Environment{}, auth.ErrUnauthorized
	}

	if !(canConfigureSensitiveEnvironments || canManageEnvironments) && env.Sensitive {
		return cdb.Environment{}, cdb.ErrEnvNotFound
	}

	return env, retrievalErr
}

func (svc *Environments) GetEnvironmentByName(ctx context.Context, actor auth.User, name string) (cdb.Environment, error) {
	env, err := svc.repo.GetEnvironmentByName(ctx, name)
	return svc.singleRetrievalPermissionChecks(
		ctx,
		actor,
		env,
		err,
	)
}

func (svc *Environments) GetEnvironmentByID(ctx context.Context, actor auth.User, id int) (cdb.Environment, error) {
	env, err := svc.repo.GetEnvironment(ctx, id)
	return svc.singleRetrievalPermissionChecks(
		ctx,
		actor,
		env,
		err,
	)
}

func (svc *Environments) ListEnvironments(ctx context.Context, actor auth.User) ([]cdb.Environment, error) {
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

func getChildren(parent cdb.Environment, environments []cdb.Environment) []cdb.EnvTree {
	children := []cdb.EnvTree{}

	for _, env := range environments {
		isChild := env.PromotesToID != nil && *env.PromotesToID == parent.ID
		if isChild {
			children = append(children, cdb.EnvTree{
				Env:      env,
				Children: getChildren(env, environments),
			})
		}
	}

	return children
}

func (svc *Environments) GetEnvironmentTree(ctx context.Context, actor auth.User) ([]cdb.EnvTree, error) {
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

	trees := []cdb.EnvTree{}
	for _, env := range environs {
		if env.PromotesToID == nil {
			trees = append(trees, cdb.EnvTree{
				Env:      env,
				Children: getChildren(env, environs),
			})
		}
	}

	return trees, nil
}
