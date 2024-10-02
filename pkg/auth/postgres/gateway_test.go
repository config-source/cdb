package postgres_test

import (
	"testing"

	pg "github.com/config-source/cdb/pkg/auth/postgres"
	"github.com/config-source/cdb/pkg/postgresutils"
	"github.com/rs/zerolog"
)

func initTestDB(t *testing.T) *pg.Gateway {
	t.Helper()

	pool := postgresutils.InitTestDB(t)
	repo := pg.NewGateway(
		zerolog.New(nil).Level(zerolog.Disabled),
		pool,
	)

	return repo
}
