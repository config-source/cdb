package auth

import (
	"context"
)

// TestGateway is an in-memory Authn/z Gateway used for tests only.
type TestGateway struct {
	IsHealthy           bool
	Error               error
	DenyPermissionCheck bool

	Users     map[UserID]User
	EmailToID map[string]UserID
}

func NewTestGateway() *TestGateway {
	return &TestGateway{
		IsHealthy:           true,
		DenyPermissionCheck: false,
		Error:               nil,
		Users:               make(map[UserID]User),
		EmailToID:           make(map[string]UserID),
	}
}

func NewTestServiceWithGateway(testgw *TestGateway) *UserService {
	return NewUserService(
		testgw,
		testgw,
		&TokenRegistry{},
		true,
		"Operator",
	)
}

func NewTestService() *UserService {
	testgw := NewTestGateway()
	return NewTestServiceWithGateway(testgw)
}

func (tr *TestGateway) Healthy(ctx context.Context) bool {
	return tr.IsHealthy
}

// AuthenticationGateway

func (tg *TestGateway) Register(ctx context.Context, email, password string) (User, error) {
	if _, ok := tg.EmailToID[email]; ok {
		return User{}, ErrEmailInUse
	}

	id := UserID(len(tg.Users) + 1)
	u := User{
		ID:       id,
		Email:    email,
		Password: password,
	}

	tg.Users[id] = u
	tg.EmailToID[email] = id
	return u, tg.Error
}

func (tg *TestGateway) Login(ctx context.Context, email, password string) (User, error) {
	userID, ok := tg.EmailToID[email]
	if !ok {
		return User{}, ErrUserNotFound
	}

	user, ok := tg.Users[userID]
	if !ok {
		return User{}, ErrUserNotFound
	}

	if user.Password == password {
		return user, tg.Error
	} else if tg.Error != nil {
		return User{}, tg.Error
	}

	return User{}, ErrInvalidPassword
}

func (tg *TestGateway) CreateUser(ctx context.Context, newUser User) (User, error) {
	return tg.Register(ctx, newUser.Email, newUser.Password)
}

func (tg *TestGateway) GetUser(ctx context.Context, userID UserID) (User, error) {
	user, ok := tg.Users[userID]
	if !ok {
		return User{}, ErrUserNotFound
	}

	return user, tg.Error
}

func (tg *TestGateway) DeleteUser(ctx context.Context, userID UserID) error {
	delete(tg.Users, userID)
	return tg.Error
}

func (tg *TestGateway) ListUsers(ctx context.Context) ([]User, error) {
	users := make([]User, len(tg.Users))
	for id, user := range tg.Users {
		users[id-1] = user
	}

	return users, tg.Error
}

// AuthorizationGateway

func (tg *TestGateway) HasPermission(
	ctx context.Context,
	actor User,
	permission Permission,
	additionalPermissions ...Permission,
) (bool, error) {
	return !tg.DenyPermissionCheck, tg.Error
}

func (tg *TestGateway) CreateRole(ctx context.Context, actor User, role string, permissions []Permission) error {
	return tg.Error
}

func (tg *TestGateway) AddPermissionsToRole(ctx context.Context, actor User, role string, permissions []Permission) error {
	return tg.Error
}

func (tg *TestGateway) RemovePermissionsFromRole(ctx context.Context, actor User, role string, permissions []Permission) error {
	return tg.Error
}

func (tg *TestGateway) GetPermissionsForRole(ctx context.Context, actor User, role string) ([]Permission, error) {
	return []Permission{PermissionConfigureEnvironments, PermissionConfigureSensitiveEnvironments}, tg.Error
}

func (tg *TestGateway) GetRolesForUser(ctx context.Context, actor User, user User) ([]string, error) {
	return []string{"test"}, tg.Error
}

func (tg *TestGateway) AssignRoleToUser(ctx context.Context, actor User, user User, role string) error {
	return tg.Error
}

func (tg *TestGateway) AssignRoleToUserNoAuth(ctx context.Context, user User, role string) error {
	return tg.Error
}

func (tg *TestGateway) RemoveRoleFromUser(ctx context.Context, actor User, user User, role string) error {
	return tg.Error
}
