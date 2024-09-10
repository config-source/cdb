package postgresutils

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

func Rollback(ctx context.Context, txn pgx.Tx, log zerolog.Logger) {
	err := txn.Rollback(ctx)
	if err != nil {
		log.Err(err).Msg("failed to rollback transaction")
	}
}
