package postgres_test

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"testing"

	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/auth/postgres"
)

var allPermissions = []auth.Permission{
	auth.PermissionConfigureEnvironments,
	auth.PermissionConfigureSensitiveEnvironments,
	auth.PermissionManageConfigKeys,
	auth.PermissionManageEnvironments,
	auth.PermissionManageUsers,
	auth.PermissionManageRoles,
}

func userFixture(t *testing.T, gw *postgres.Gateway, email string) auth.User {
	user, err := gw.Register(context.Background(), email, "test")
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func userFixtureWithRole(t *testing.T, gw *postgres.Gateway, email string, role string) auth.User {
	t.Helper()

	user := userFixture(t, gw, email)
	err := gw.AssignRoleToUserNoAuth(context.Background(), user, role)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func adminFixture(t *testing.T, gw *postgres.Gateway) auth.User {
	t.Helper()
	return userFixtureWithRole(t, gw, "admin@example.com", "Administrator")
}

func operatorFixture(t *testing.T, gw *postgres.Gateway) auth.User {
	t.Helper()
	return userFixtureWithRole(t, gw, "operator@example.com", "Operator")
}

func TestAdminHasAllPermission(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	admin := adminFixture(t, gateway)
	ctx := context.Background()

	// Admin should have all permissions
	for _, perm := range allPermissions {
		hasPerm, err := gateway.HasPermission(ctx, admin, perm)
		if err != nil {
			t.Error(err)
		}

		if !hasPerm {
			t.Errorf("Expected Administrator to have permission: %s", perm)
		}
	}
}

func TestOperatorHasConfigurePermissions(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	operator := operatorFixture(t, gateway)
	ctx := context.Background()

	// Operator should have all configuration permissions
	for _, perm := range []auth.Permission{
		auth.PermissionConfigureEnvironments,
		auth.PermissionConfigureSensitiveEnvironments,
	} {
		hasPerm, err := gateway.HasPermission(ctx, operator, perm)
		if err != nil {
			t.Error(err)
		}

		if !hasPerm {
			t.Errorf("Expected Administrator to have permission: %s", perm)
		}
	}
}

func TestAssignRoleToUser(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	admin := adminFixture(t, gateway)
	unassignedUser := userFixture(t, gateway, "user@example.com")
	ctx := context.Background()

	err := gateway.AssignRoleToUser(ctx, admin, unassignedUser, "Operator")
	if err != nil {
		t.Fatal(err)
	}

	roles, err := gateway.GetRolesForUser(ctx, admin, unassignedUser)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(roles, []string{"Operator"}) {
		t.Errorf("Expected %s roles but got: %s", []string{"Operator"}, roles)
	}
}

func TestAssignNonExistentRoleToUser(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	admin := adminFixture(t, gateway)
	unassignedUser := userFixture(t, gateway, "user@example.com")
	ctx := context.Background()

	err := gateway.AssignRoleToUser(ctx, admin, unassignedUser, "NonExistent")
	if err == nil {
		t.Error("Expected an error but got nil!")
	}
}

func TestAssignRoleToUserRequiresManageUserPermission(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	operator := operatorFixture(t, gateway)
	unassignedUser := userFixture(t, gateway, "user@example.com")
	ctx := context.Background()

	err := gateway.AssignRoleToUser(ctx, operator, unassignedUser, "Operator")
	if !errors.Is(err, auth.ErrUnauthorized) {
		t.Errorf("Expected %s got: %s", auth.ErrUnauthorized, err)
	}
}

func TestRemoveRoleFromUser(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	admin := adminFixture(t, gateway)
	unassignedUser := userFixture(t, gateway, "user@example.com")
	ctx := context.Background()

	err := gateway.AssignRoleToUser(ctx, admin, unassignedUser, "Operator")
	if err != nil {
		t.Fatal(err)
	}

	err = gateway.RemoveRoleFromUser(ctx, admin, unassignedUser, "Operator")
	if err != nil {
		t.Fatal(err)
	}

	roles, err := gateway.GetRolesForUser(ctx, admin, unassignedUser)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(roles, []string{}) {
		t.Errorf("Expected %s roles but got: %s", []string{}, roles)
	}
}

func TestRemoveNonExistentRoleFromUserReturnsError(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	admin := adminFixture(t, gateway)
	unassignedUser := userFixture(t, gateway, "user@example.com")
	ctx := context.Background()

	err := gateway.RemoveRoleFromUser(ctx, admin, unassignedUser, "Operator")
	if err == nil {
		t.Error("Expected an error but got nil!")
	}
}

func TestGetPermissionsForAdministratorRole(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	admin := adminFixture(t, gateway)
	ctx := context.Background()

	perms, err := gateway.GetPermissionsForRole(ctx, admin, "Administrator")
	if err != nil {
		t.Fatal(err)
	}

	permStrings := make([]string, len(perms))
	for i, perm := range permStrings {
		permStrings[i] = string(perm)
	}

	allPermStrings := make([]string, len(allPermissions))
	for i, perm := range allPermStrings {
		allPermStrings[i] = string(perm)
	}

	if !reflect.DeepEqual(
		sort.StringSlice(permStrings),
		sort.StringSlice(allPermStrings),
	) {
		t.Errorf("Expected %s permissions but got: %s", allPermissions, perms)
	}
}

func TestGetPermissionsForOperatorRole(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	admin := adminFixture(t, gateway)
	ctx := context.Background()

	perms, err := gateway.GetPermissionsForRole(ctx, admin, "Operator")
	if err != nil {
		t.Fatal(err)
	}

	expectedPerms := []auth.Permission{
		auth.PermissionConfigureEnvironments,
		auth.PermissionConfigureSensitiveEnvironments,
	}
	if !reflect.DeepEqual(perms, expectedPerms) {
		t.Errorf("Expected %s permissions but got: %s", expectedPerms, perms)
	}
}

func TestRemovePermissionsFromRole(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	admin := adminFixture(t, gateway)
	ctx := context.Background()

	err := gateway.RemovePermissionsFromRole(
		ctx,
		admin,
		"Operator",
		[]auth.Permission{
			auth.PermissionConfigureSensitiveEnvironments,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	perms, err := gateway.GetPermissionsForRole(ctx, admin, "Operator")
	if err != nil {
		t.Fatal(err)
	}

	expectedPerms := []auth.Permission{
		auth.PermissionConfigureEnvironments,
	}
	if !reflect.DeepEqual(perms, expectedPerms) {
		t.Errorf("Expected %s permissions but got: %s", expectedPerms, perms)
	}
}

func TestAddPermissionsToRole(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	admin := adminFixture(t, gateway)
	ctx := context.Background()

	err := gateway.RemovePermissionsFromRole(
		ctx,
		admin,
		"Operator",
		[]auth.Permission{
			auth.PermissionConfigureSensitiveEnvironments,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = gateway.AddPermissionsToRole(
		ctx,
		admin,
		"Operator",
		[]auth.Permission{
			auth.PermissionConfigureSensitiveEnvironments,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	perms, err := gateway.GetPermissionsForRole(ctx, admin, "Operator")
	if err != nil {
		t.Fatal(err)
	}

	expectedPerms := []auth.Permission{
		auth.PermissionConfigureEnvironments,
		auth.PermissionConfigureSensitiveEnvironments,
	}
	if !reflect.DeepEqual(perms, expectedPerms) {
		t.Errorf("Expected %s permissions but got: %s", expectedPerms, perms)
	}
}

func TestCreateRole(t *testing.T) {
	gateway, tr := initTestDB(t)
	defer tr.Cleanup()

	admin := adminFixture(t, gateway)
	unassignedUser := userFixture(t, gateway, "user@example.com")
	newRoleName := "Minion"
	ctx := context.Background()

	err := gateway.CreateRole(ctx, admin, newRoleName, []auth.Permission{})
	if err != nil {
		t.Fatal(err)
	}

	err = gateway.AssignRoleToUser(ctx, admin, unassignedUser, newRoleName)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := gateway.GetRolesForUser(ctx, admin, unassignedUser)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(roles, []string{newRoleName}) {
		t.Errorf("Expected %s roles but got: %s", []string{newRoleName}, roles)
	}
}
