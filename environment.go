package cdb

import (
	"context"
	"errors"
	"fmt"
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

func (e Environment) String() string {
	promotesToID := 0
	if e.PromotesToID != nil {
		promotesToID = *e.PromotesToID
	}

	return fmt.Sprintf(
		"Environment(id=%d, name=%s, promotes_to=%d)",
		e.ID,
		e.Name,
		promotesToID,
	)
}

type EnvironmentRepository interface {
	CreateEnvironment(context.Context, Environment) (Environment, error)

	GetEnvironment(ctx context.Context, id int) (Environment, error)
	GetEnvironmentByName(ctx context.Context, name string) (Environment, error)

	ListEnvironments(ctx context.Context) ([]Environment, error)
}

type EnvTree struct {
	Env      Environment
	Children []EnvTree
}
