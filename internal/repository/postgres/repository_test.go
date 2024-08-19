package postgres_test

import (
	"context"
	"testing"

	"github.com/config-source/cdb/internal/postgresutils"
	pg "github.com/config-source/cdb/internal/repository/postgres"
	"github.com/rs/zerolog"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func initTestDB(t *testing.T) (*pg.Repository, *postgresutils.TestRepository) {
	t.Helper()

	tr := postgresutils.InitTestDB(t)

	repo, err := pg.NewRepository(
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
