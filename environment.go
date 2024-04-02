package cdb

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEnvNotFound = errors.New("environment not found")
)

type Environment struct {
	ID int `db:"id"`

	Name         string `db:"name"`
	PromotesToID *int   `db:"promotes_to_id"`

	CreatedAt time.Time `db:"created_at"`
}

type EnvironmentRepository interface {
	CreateEnvironment(context.Context, Environment) (Environment, error)

	GetEnvironment(ctx context.Context, id int) (Environment, error)
	GetEnvironmentByName(ctx context.Context, name string) (Environment, error)
}
