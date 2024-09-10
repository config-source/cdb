package postgres

import (
	"context"

	"github.com/config-source/cdb/postgresutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Repository struct {
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func NewRepository(ctx context.Context, log zerolog.Logger, connString string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &Repository{
		pool: pool,
		log:  log,
	}, nil
}

// Raw returns the raw pgxpool.Pool in use by this repository. Should only be
// used for testing purposes.
func (r *Repository) Raw() *pgxpool.Pool {
	return r.pool
}

// Healthy returns a boolean indicating that the connection pool is working and
// queries can be run.
func (r *Repository) Healthy(ctx context.Context) bool {
	return postgresutils.HealthCheck(ctx, r.pool, r.log)
}
