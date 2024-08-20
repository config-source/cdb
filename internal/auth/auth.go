package auth

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrEmailInUse      = errors.New("email is already in use")
	ErrUnauthorized    = errors.New("you do not have permission to perform that action")
	ErrUnauthenticated = errors.New("you do not have permission to perform that action")
)

// UserID is a custom type for user IDs to force some validation and to ease
// changing the type later as we support additional gateways.
type UserID int

// User is an internal user of the system.
type User struct {
	ID UserID `db:"id"`

	Email string `db:"email"`
	// Completely ignored in JSON to avoid any accidental leakages. Incoming
	// requests have a dedicated write-only type that's stored in the api
	// package.
	Password string `db:"password" json:"-"`
}

func (u User) String() string {
	return fmt.Sprintf("User(email=%s)", u.Email)
}

// AuthenticationGateway must be implemented by any source of authentication in
// CDB.
//
// Currently only internal DB-based authentication is supported. But this
// interface is here so we can implement it for other backends later like LDAP
// etc.
type AuthenticationGateway interface {
	Register(ctx context.Context, email, password string) (User, error)
	Login(ctx context.Context, email, password string) (User, error)

	CreateUser(ctx context.Context, newUser User) (User, error)
	GetUser(ctx context.Context, userID UserID) (User, error)
	DeleteUser(ctx context.Context, userID UserID) error
	ListUsers(ctx context.Context) ([]User, error)

	Healthy(context.Context) bool
}

// AuthorizationGateway must be implemented by any source of authorization in
// CDB.
//
// Currently only internal DB-based authorization is supported. But this
// interface is here so we can implement it for other backends later like LDAP
// etc.
type AuthorizationGateway interface {
	// TODO: should this just return an error or just a bool?
	HasPermission(ctx context.Context, actor User, permission Permission) (bool, error)

	CreateRole(ctx context.Context, actor User, role string, permissions []Permission) error
	AddPermissionsToRole(ctx context.Context, actor User, role string, permissions []Permission) error
	RemovePermissionsFromRole(ctx context.Context, actor User, role string, permissions []Permission) error
	GetPermissionsForRole(ctx context.Context, actor User, role string) ([]Permission, error)

	GetRolesForUser(ctx context.Context, actor User, user User) ([]string, error)
	AssignRoleToUser(ctx context.Context, actor User, user User, role string) error
	RemoveRoleFromUser(ctx context.Context, actor User, user User, role string) error

	Healthy(context.Context) bool
}
