package auth

import (
	"context"
	"errors"
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

// TokenSet is a set of JWT tokens for use as Authentication and Authorisation.
type TokenSet struct {
	IDToken      string
	AccessToken  string
	RefreshToken string
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

	CreateUser(ctx context.Context, actor User, newUser User) (User, error)
	GetUser(ctx context.Context, actor User, userID UserID) (User, error)
	DeleteUser(ctx context.Context, actor User, userID UserID) error
	ListUsers(ctx context.Context, actor User) ([]User, error)

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
	HasPermission(ctx context.Context, actor User, permission string) (bool, error)

	CreateRole(ctx context.Context, actor User, role string, permissions []string) error
	AddPermissionsToRole(ctx context.Context, actor User, role string, permissions []string) error
	RemovePermissionsFromRole(ctx context.Context, actor User, role string, permissions []string) error
	GetPermissionsForRole(ctx context.Context, actor User, role string) ([]string, error)

	GetRolesForUser(ctx context.Context, actor User, user User) ([]string, error)
	AssignRoleToUser(ctx context.Context, actor User, user User, role string) error
	RemoveRoleFromUser(ctx context.Context, actor User, user User, role string) error

	Healthy(context.Context) bool
}
