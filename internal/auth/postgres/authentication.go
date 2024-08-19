package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/config-source/cdb/internal/auth"
	"github.com/config-source/cdb/internal/postgresutils"
	"golang.org/x/crypto/bcrypt"
)

//go:embed queries/authentication/create_user.sql
var createUserSql string

//go:embed queries/authentication/get_user_by_email.sql
var getUserByEmailSql string

func (g *Gateway) Register(ctx context.Context, email, password string) (auth.User, error) {
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

func (g *Gateway) CreateUser(ctx context.Context, actor auth.User, newUser auth.User) (auth.User, error) {
	// TODO: need to implement permissions first.
	return auth.User{}, nil
}

func (g *Gateway) GetUser(ctx context.Context, actor auth.User, userID auth.UserID) (auth.User, error) {
	// TODO: need to implement permissions first.
	return auth.User{}, nil
}

func (g *Gateway) DeleteUser(ctx context.Context, actor auth.User, userID auth.UserID) error {
	// TODO: need to implement permissions first.
	return nil
}

func (g *Gateway) ListUsers(ctx context.Context, actor auth.User) ([]auth.User, error) {
	// TODO: need to implement permissions first.
	return []auth.User{}, nil
}
