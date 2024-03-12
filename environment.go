package cdb

import uuid "github.com/gofrs/uuid/v5"

type Environment struct {
	ID uuid.UUID `db:"id"`

	Name       string     `db:"name"`
	PromotesTo *uuid.UUID `db:"promotes_to"`
}

type EnvironmentRepository interface {
	CreateEnvironment(Environment) (Environment, error)

	GetEnvironment(id uuid.UUID) (Environment, error)
	GetEnvironmentByName(name string) (Environment, error)
}
