package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/config-source/cdb/pkg/auth"
	"github.com/config-source/cdb/pkg/postgresutils"
	"golang.org/x/crypto/bcrypt"
)

//go:embed queries/authentication/create_user.sql
var createUserSql string

//go:embed queries/authentication/get_user_by_email.sql
var getUserByEmailSql string

//go:embed queries/authentication/get_user_by_id.sql
var getUserByIDSql string

//go:embed queries/authentication/get_users.sql
var getUsersSql string

//go:embed queries/authentication/delete_user.sql
var deleteUserSql string

func (g *Gateway) Register(ctx context.Context, email, password string) (auth.User, error) {
	// TODO: enforce some password strength
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return auth.User{}, err
	}

	user, err := postgresutils.GetOne[auth.User](g.pool, ctx, createUserSql, email, hashedPw)
	if err != nil && postgresutils.IsUniqueConstraintErr(err) {
		return auth.User{}, fmt.Errorf("%w: %s", auth.ErrEmailInUse, email)
	}

	return user, err
}

func (g *Gateway) Login(ctx context.Context, email, password string) (auth.User, error) {
	user, err := postgresutils.GetOne[auth.User](g.pool, ctx, getUserByEmailSql, email)
	if err != nil {
		return auth.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return auth.User{}, auth.ErrInvalidPassword
	}

	return user, nil
}

func (g *Gateway) CreateUser(ctx context.Context, newUser auth.User) (auth.User, error) {
	// TODO: in this scenario we want to notify the new user. I don't really
	// think that should live here though probably in a wrapping Service type.
	return g.Register(ctx, newUser.Email, newUser.Password)
}

func (g *Gateway) GetUser(ctx context.Context, userID auth.UserID) (auth.User, error) {
	// TODO: return user not found as appropriate
	return postgresutils.GetOne[auth.User](g.pool, ctx, getUserByIDSql, userID)
}

func (g *Gateway) DeleteUser(ctx context.Context, userID auth.UserID) error {
	_, err := g.pool.Exec(ctx, deleteUserSql, userID)
	return err
}

func (g *Gateway) ListUsers(ctx context.Context) ([]auth.User, error) {
	return postgresutils.GetAll[auth.User](g.pool, ctx, getUsersSql)
}

// TODO: update user
