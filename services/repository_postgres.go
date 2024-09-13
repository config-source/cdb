package services

import (
	"context"
	_ "embed"
	"errors"

	"github.com/config-source/cdb/postgresutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func NewRepository(log zerolog.Logger, pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		log:  log,
		pool: pool,
	}
}

//go:embed queries/create_service.sql
var createServiceSql string

//go:embed queries/get_service_by_id.sql
var getServiceByIDSql string

//go:embed queries/get_service_by_name.sql
var getServiceByNameSql string

//go:embed queries/list_services.sql
var listServicesSql string

func (r *PostgresRepository) CreateService(ctx context.Context, svc Service) (Service, error) {
	return postgresutils.GetOne[Service](r.pool, ctx, createServiceSql, svc.Name)
}

func (r *PostgresRepository) GetService(ctx context.Context, id int) (Service, error) {
	svc, err := postgresutils.GetOne[Service](r.pool, ctx, getServiceByIDSql, id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return svc, ErrNotFound
	}

	return svc, err
}

func (r *PostgresRepository) GetServiceByName(ctx context.Context, name string) (Service, error) {
	svc, err := postgresutils.GetOne[Service](r.pool, ctx, getServiceByNameSql, name)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return svc, ErrNotFound
	}

	return svc, err
}

func (r *PostgresRepository) ListServices(ctx context.Context, includeSensitive bool) ([]Service, error) {
	return postgresutils.GetAll[Service](r.pool, ctx, listServicesSql)
}
