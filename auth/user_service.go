package auth

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrPublicRegisterDisabled = errors.New("public registration is disabled")
)

type UserService struct {
	authn    AuthenticationGateway
	authz    AuthorizationGateway
	registry *TokenRegistry

	publicRegisterAllowed bool
	defaultRegisterRole   string
}

func NewUserService(
	authn AuthenticationGateway,
	authz AuthorizationGateway,
	registry *TokenRegistry,
	allowPublicRegistration bool,
	defaultRegisterRole string,
) *UserService {
	return &UserService{
		authn:                 authn,
		authz:                 authz,
		registry:              registry,
		publicRegisterAllowed: allowPublicRegistration,
		defaultRegisterRole:   defaultRegisterRole,
	}
}

func (us *UserService) Healthy(ctx context.Context) bool {
	return us.authn.Healthy(ctx)
}

func (us *UserService) Login(ctx context.Context, email, password string) (User, error) {
	return us.authn.Login(ctx, email, password)
}

func (us *UserService) Logout(ctx context.Context, refreshToken string) error {
	return us.registry.Revoke(ctx, refreshToken)
}

func (us *UserService) Register(ctx context.Context, email, password string) (User, error) {
	if us.publicRegisterAllowed {
		user, err := us.authn.Register(ctx, email, password)
		if err != nil {
			return user, err
		}

		return user, us.authz.AssignRoleToUserNoAuth(ctx, user, us.defaultRegisterRole)
	}

	return User{}, ErrPublicRegisterDisabled
}

func (us *UserService) CreateUser(ctx context.Context, actor User, newUser User) (User, error) {
	isAuthorized, err := us.authz.HasPermission(ctx, actor, PermissionManageUsers)
	if err != nil {
		return User{}, err
	}

	if !isAuthorized {
		return User{}, fmt.Errorf("%w: not permitted to manage users", ErrUnauthorized)
	}

	return us.authn.CreateUser(ctx, newUser)
}

func (us *UserService) GetUser(ctx context.Context, actor User, userID UserID) (User, error) {
	isAuthorized, err := us.authz.HasPermission(ctx, actor, PermissionManageUsers)
	if err != nil {
		return User{}, err
	}

	if !isAuthorized {
		return User{}, fmt.Errorf("%w: not permitted to manage users", ErrUnauthorized)
	}

	return us.authn.GetUser(ctx, userID)
}

func (us *UserService) DeleteUser(ctx context.Context, actor User, userID UserID) error {
	isAuthorized, err := us.authz.HasPermission(ctx, actor, PermissionManageUsers)
	if err != nil {
		return err
	}

	if !isAuthorized {
		return fmt.Errorf("%w: not permitted to manage users", ErrUnauthorized)
	}

	return us.authn.DeleteUser(ctx, userID)
}

func (us *UserService) ListUsers(ctx context.Context, actor User) ([]User, error) {
	isAuthorized, err := us.authz.HasPermission(ctx, actor, PermissionManageUsers)
	if err != nil {
		return nil, err
	}

	if !isAuthorized {
		return nil, fmt.Errorf("%w: not permitted to manage users", ErrUnauthorized)
	}

	return us.authn.ListUsers(ctx)
}

func (us *UserService) IssueAPIToken(ctx context.Context, signingKey []byte, actor User) (APIToken, error) {
	return us.registry.IssueAPIToken(ctx, signingKey, actor)
}

func (us *UserService) ListAPITokens(ctx context.Context, actor User) ([]APIToken, error) {
	return us.registry.ListAPITokens(ctx, actor)
}
