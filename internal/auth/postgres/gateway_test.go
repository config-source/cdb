package postgres_test

import (
	"context"
	"testing"

	pg "github.com/config-source/cdb/internal/auth/postgres"
	"github.com/config-source/cdb/internal/postgresutils"
	"github.com/rs/zerolog"
)

func initTestDB(t *testing.T) (*pg.Gateway, *postgresutils.TestRepository) {
	t.Helper()

	tr := postgresutils.InitTestDB(t)

	repo, err := pg.NewGateway(
		context.Background(),
		zerolog.New(nil).Level(zerolog.Disabled),
		tr.TestDBURL,
	)
	if err != nil {
		t.Fatal(err)
	}
	tr.Repo = repo

	return repo, tr
}
