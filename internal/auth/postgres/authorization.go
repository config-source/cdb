package postgres

import (
	"context"

	"github.com/config-source/cdb/internal/auth"
)

func (g *Gateway) HasPermission(ctx context.Context, actor auth.User, permission string) (bool, error) {
	return false, nil
}

func (g *Gateway) CreateRole(ctx context.Context, actor auth.User, role string, permissions []string) error {
	return nil
}

func (g *Gateway) AddPermissionsToRole(ctx context.Context, actor auth.User, role string, permissions []string) error {
	return nil
}

func (g *Gateway) RemovePermissionsFromRole(ctx context.Context, actor auth.User, role string, permissions []string) error {
	return nil
}

func (g *Gateway) GetPermissionsForRole(ctx context.Context, actor auth.User, role string) ([]string, error) {
	return []string{}, nil
}

func (g *Gateway) GetRolesForauthUser(ctx context.Context, actor auth.User, user auth.User) ([]string, error) {
	return []string{}, nil
}

func (g *Gateway) AssignRoleToauthUser(ctx context.Context, actor auth.User, user auth.User, role string) error {
	return nil
}

func (g *Gateway) RemoveRoleFromauthUser(ctx context.Context, actor auth.User, user auth.User, role string) error {
	return nil
}
