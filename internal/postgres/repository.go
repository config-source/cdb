package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(ctx context.Context, connString string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &Repository{
		pool: pool,
	}, nil
}

// Raw returns the raw pgxpool.Pool in use by this repository. Should only be
// used for development purposes.
func (r *Repository) Raw() *pgxpool.Pool {
	return r.pool
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
