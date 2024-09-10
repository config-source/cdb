package repository

import (
	"context"

	"github.com/config-source/cdb"
)

type ModelRepository interface {
	cdb.EnvironmentRepository
	cdb.ConfigValueRepository
	cdb.ConfigKeyRepository
	Healthy(context.Context) bool
}
