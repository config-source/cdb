package postgres_test

import (
	"testing"

	pg "github.com/config-source/cdb/auth/postgres"
	"github.com/config-source/cdb/postgresutils"
	"github.com/rs/zerolog"
)

func initTestDB(t *testing.T) (*pg.Gateway, *postgresutils.TestRepository) {
	t.Helper()

	tr, pool := postgresutils.InitTestDB(t)
	repo := pg.NewGateway(
		zerolog.New(nil).Level(zerolog.Disabled),
		pool,
	)

	return repo, tr
}
