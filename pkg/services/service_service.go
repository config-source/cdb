package services

import (
	"context"

	"github.com/config-source/cdb/pkg/auth"
)

type ServiceService struct {
	auth auth.AuthorizationGateway
	repo *Repository
}

func NewServiceService(repo *Repository, auth auth.AuthorizationGateway) *ServiceService {
	return &ServiceService{
		auth: auth,
		repo: repo,
	}
}

func (s *ServiceService) CreateService(ctx context.Context, actor auth.User, svc Service) (Service, error) {
	canManageEnvironments, err := s.auth.HasPermission(ctx, actor, auth.PermissionManageEnvironments)
	if err != nil {
		return Service{}, err
	}

	if !canManageEnvironments {
		return Service{}, auth.ErrUnauthorized
	}

	return s.repo.CreateService(ctx, svc)
}

func (s *ServiceService) viewPermissionCheck(ctx context.Context, actor auth.User) error {
	canViewServices, err := s.auth.HasPermission(
		ctx,
		actor,
		auth.PermissionManageEnvironments,
		auth.PermissionConfigureEnvironments,
		auth.PermissionConfigureSensitiveEnvironments,
	)
	if err != nil {
		return err
	}

	if !canViewServices {
		return auth.ErrUnauthorized
	}

	return nil
}

func (s *ServiceService) GetServiceByName(ctx context.Context, actor auth.User, name string) (Service, error) {
	if err := s.viewPermissionCheck(ctx, actor); err != nil {
		return Service{}, err
	}

	return s.repo.GetServiceByName(ctx, name)
}

func (s *ServiceService) GetServiceByID(ctx context.Context, actor auth.User, id int) (Service, error) {
	if err := s.viewPermissionCheck(ctx, actor); err != nil {
		return Service{}, err
	}

	return s.repo.GetService(ctx, id)
}

func (s *ServiceService) ListServices(ctx context.Context, actor auth.User) ([]Service, error) {
	if err := s.viewPermissionCheck(ctx, actor); err != nil {
		return nil, err
	}

	return s.repo.ListServices(ctx, true)
}
