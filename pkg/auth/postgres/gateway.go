package postgres

import (
	"context"

	"github.com/config-source/cdb/pkg/postgresutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// Gateway implements auth.AuthenticationGateway and auth.AuthorizationGateway
// using Postgres as a storage backend.
type Gateway struct {
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func NewGateway(log zerolog.Logger, pool *pgxpool.Pool) *Gateway {
	return &Gateway{
		pool: pool,
		log:  log,
	}
}

// Raw returns the raw pgxpool.Pool in use by this repository. Should only be
// used for testing purposes.
func (g *Gateway) Raw() *pgxpool.Pool {
	return g.pool
}

// Healthy returns a boolean indicating that the connection pool is working and
// queries can be run.
func (g *Gateway) Healthy(ctx context.Context) bool {
	return postgresutils.HealthCheck(ctx, g.pool, g.log)
}
