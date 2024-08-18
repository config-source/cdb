package postgresutils

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// HealthCheck confirms that the given pool can connect to Postgres
func HealthCheck(ctx context.Context, pool *pgxpool.Pool, logger zerolog.Logger) bool {
	var healthy int
	var log *zerolog.Event

	err := pool.QueryRow(ctx, "SELECT 1").Scan(&healthy)
	if err != nil {
		log = logger.Err(err)
	} else if healthy != 1 {
		log = logger.Error()
	} else {
		log = logger.Info()
	}

	log.
		Int("healthy", healthy).
		Msg("Postgres health check")

	return healthy == 1
}

// boilerplate reducing utilities

// GetOne runs the given query and serializes the returned row into T
//
// If more than a single row matches the given query an error is returned.
func GetOne[T any](pool *pgxpool.Pool, ctx context.Context, sql string, args ...interface{}) (T, error) {
	rows, err := pool.Query(ctx, sql, args...)
	if err != nil {
		var def T
		return def, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[T])
}

// GetAll runs the given query and serializes the returned rows into a slice of T
func GetAll[T any](pool *pgxpool.Pool, ctx context.Context, sql string, args ...interface{}) ([]T, error) {
	rows, err := pool.Query(ctx, sql, args...)
	if err != nil {
		var def []T
		return def, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

// IsUniqueConstraintErr returns a boolean indicating if the given error is
// a unique constraint violation so callers can respond accordingly.
func IsUniqueConstraintErr(err error) bool {
	return strings.Contains(err.Error(), "unique constraint")
}
