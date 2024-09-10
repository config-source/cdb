package postgres

import (
	"context"
	_ "embed"
	"errors"

	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/postgresutils"
	"github.com/jackc/pgx/v5"
)

//go:embed queries/authorization/has_permission.sql
var hasPermissionSql string

//go:embed queries/authorization/assign_role_to_user.sql
var assignRoleToUserSql string

//go:embed queries/authorization/remove_role_from_user.sql
var removeRoleFromUserSql string

//go:embed queries/authorization/get_roles_for_user.sql
var getRolesForUserSql string

//go:embed queries/authorization/get_permissions_for_role.sql
var getPermissionsForRoleSql string

//go:embed queries/authorization/create_role.sql
var createRoleSql string

//go:embed queries/authorization/assign_permission_to_role.sql
var assignPermissionToRoleSql string

//go:embed queries/authorization/remove_permission_from_role.sql
var removePermissionFromRoleSql string

func (g *Gateway) HasPermission(ctx context.Context, actor auth.User, permission auth.Permission) (bool, error) {
	var rowCount int
	err := g.pool.QueryRow(ctx, hasPermissionSql, actor.ID, permission).Scan(&rowCount)
	return rowCount > 0, err
}

func (g *Gateway) CreateRole(ctx context.Context, actor auth.User, role string, permissions []auth.Permission) error {
	if isAuthorized, _ := g.HasPermission(ctx, actor, auth.PermissionManageRoles); !isAuthorized {
		g.log.Warn().
			Interface("actorID", actor.ID).
			Str("role", role).
			Interface("permissions", permissions).
			Bool("denied", true).
			Bool("audit", true).
			Msg("actor attempted to create a role")

		return auth.ErrUnauthorized
	}

	_, err := g.pool.Exec(ctx, createRoleSql, role)
	if err != nil {
		return err
	}

	g.log.Info().
		Interface("actorID", actor.ID).
		Str("role", role).
		Interface("permissions", permissions).
		Bool("denied", false).
		Bool("audit", true).
		Msg("actor created a role")

	return g.AddPermissionsToRole(ctx, actor, role, permissions)
}

func (g *Gateway) AddPermissionsToRole(ctx context.Context, actor auth.User, role string, permissions []auth.Permission) error {
	if isAuthorized, _ := g.HasPermission(ctx, actor, auth.PermissionManageRoles); !isAuthorized {
		g.log.Warn().
			Interface("actorID", actor.ID).
			Str("role", role).
			Interface("permissions", permissions).
			Bool("denied", true).
			Bool("audit", true).
			Msg("actor attempted to add permissions to role")

		return auth.ErrUnauthorized
	}

	txn, err := g.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		postgresutils.Rollback(ctx, txn, g.log)
		return err
	}

	for _, perm := range permissions {
		_, assignErr := txn.Exec(ctx, assignPermissionToRoleSql, role, perm)
		if assignErr != nil {
			postgresutils.Rollback(ctx, txn, g.log)
			return assignErr
		}
	}

	g.log.Info().
		Interface("actorID", actor.ID).
		Str("role", role).
		Interface("permissions", permissions).
		Bool("denied", false).
		Bool("audit", true).
		Msg("actor added permissions to role")

	return txn.Commit(ctx)
}

func (g *Gateway) RemovePermissionsFromRole(ctx context.Context, actor auth.User, role string, permissions []auth.Permission) error {
	if isAuthorized, _ := g.HasPermission(ctx, actor, auth.PermissionManageRoles); !isAuthorized {
		g.log.Warn().
			Interface("actorID", actor.ID).
			Str("role", role).
			Interface("permissions", permissions).
			Bool("denied", false).
			Bool("audit", true).
			Msg("actor attempted to remove permissions from role")

		return auth.ErrUnauthorized
	}

	txn, err := g.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		postgresutils.Rollback(ctx, txn, g.log)
		return err
	}

	for _, perm := range permissions {
		_, removeErr := txn.Exec(ctx, removePermissionFromRoleSql, role, perm)
		if removeErr != nil {
			postgresutils.Rollback(ctx, txn, g.log)
			return removeErr
		}
	}

	g.log.Info().
		Interface("actorID", actor.ID).
		Str("role", role).
		Interface("permissions", permissions).
		Bool("denied", false).
		Bool("audit", true).
		Msg("actor removed permissions from role")

	return txn.Commit(ctx)
}

func (g *Gateway) GetPermissionsForRole(ctx context.Context, actor auth.User, role string) ([]auth.Permission, error) {
	if isAuthorized, _ := g.HasPermission(ctx, actor, auth.PermissionManageRoles); !isAuthorized {
		g.log.Warn().
			Interface("actorID", actor.ID).
			Str("role", role).
			Bool("denied", true).
			Bool("audit", true).
			Msg("actor attempted to get permissions for a role")

		return nil, auth.ErrUnauthorized
	}

	rows, err := g.pool.Query(ctx, getPermissionsForRoleSql, role)
	if err != nil {
		return nil, err
	}

	g.log.Info().
		Interface("actorID", actor.ID).
		Str("role", role).
		Bool("denied", false).
		Bool("audit", true).
		Msg("actor got permissions for a role")

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (auth.Permission, error) {
		var permissionName auth.Permission
		return permissionName, row.Scan(&permissionName)
	})

}

func (g *Gateway) GetRolesForUser(ctx context.Context, actor auth.User, user auth.User) ([]string, error) {
	if isAuthorized, _ := g.HasPermission(ctx, actor, auth.PermissionManageUsers); !isAuthorized {
		g.log.Warn().
			Interface("actorID", actor.ID).
			Interface("user", user.ID).
			Bool("denied", true).
			Bool("audit", true).
			Msg("actor attempted to get roles for user")

		return nil, auth.ErrUnauthorized
	}

	rows, err := g.pool.Query(ctx, getRolesForUserSql, user.ID)
	if err != nil {
		return nil, err
	}

	g.log.Info().
		Interface("actorID", actor.ID).
		Interface("user", user.ID).
		Bool("denied", false).
		Bool("audit", true).
		Msg("actor got roles for user")

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
		var roleName string
		return roleName, row.Scan(&roleName)
	})
}

func (g *Gateway) AssignRoleToUser(ctx context.Context, actor auth.User, user auth.User, role string) error {
	if isAuthorized, _ := g.HasPermission(ctx, actor, auth.PermissionManageUsers); !isAuthorized {
		g.log.Warn().
			Interface("actorID", actor.ID).
			Interface("user", user.ID).
			Str("role", role).
			Bool("denied", true).
			Bool("audit", true).
			Msg("actor attempted to assign role to user")

		return auth.ErrUnauthorized
	}

	err := g.AssignRoleToUserNoAuth(ctx, user, role)
	if err != nil {
		return err
	}

	g.log.Info().
		Interface("actorID", actor.ID).
		Interface("user", user.ID).
		Str("role", role).
		Bool("denied", false).
		Bool("audit", true).
		Msg("actor assigned role to user")

	return nil
}

func (g *Gateway) AssignRoleToUserNoAuth(ctx context.Context, user auth.User, role string) error {
	var roleID int
	err := g.pool.QueryRow(ctx, "SELECT id FROM roles WHERE name = $1", role).Scan(&roleID)
	if err != nil {
		return err
	}

	_, err = g.pool.Exec(ctx, assignRoleToUserSql, user.ID, roleID)
	return err
}

func (g *Gateway) RemoveRoleFromUser(ctx context.Context, actor auth.User, user auth.User, role string) error {
	if isAuthorized, _ := g.HasPermission(ctx, actor, auth.PermissionManageUsers); !isAuthorized {
		g.log.Warn().
			Interface("actorID", actor.ID).
			Interface("user", user.ID).
			Str("role", role).
			Bool("denied", true).
			Bool("audit", true).
			Msg("actor attempted to remove role from user")

		return auth.ErrUnauthorized
	}

	commandTag, err := g.pool.Exec(ctx, removeRoleFromUserSql, user.ID, role)
	if commandTag.RowsAffected() <= 0 {
		return errors.New("given role does not exist or was not assigned to the user")
	}

	g.log.Info().
		Interface("actorID", actor.ID).
		Interface("user", user.ID).
		Str("role", role).
		Bool("denied", false).
		Bool("audit", true).
		Msg("actor removed role from user")

	return err
}
