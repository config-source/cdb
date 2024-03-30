package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
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
// used for development purposes.
func (r *Repository) Raw() *pgxpool.Pool {
	return r.pool
}

// Healthy returns a boolean indicating that the connection pool is working and
// queries can be run.
func (r *Repository) Healthy(ctx context.Context) bool {
	var healthy int
	err := r.pool.QueryRow(ctx, "SELECT 1 FROM environments LIMIT 1").Scan(&healthy)
	if err != nil {
		r.log.Err(err)
		return false
	}

	return healthy == 1
}

// boilerplate reducing utilities

func getOne[T any](r *Repository, ctx context.Context, sql string, args ...interface{}) (T, error) {
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		var def T
		return def, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[T])
}

func getAll[T any](r *Repository, ctx context.Context, sql string, args ...interface{}) ([]T, error) {
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		var def []T
		return def, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}
