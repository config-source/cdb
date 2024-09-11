package auth

import (
	"context"
	"strconv"
	"time"

	"github.com/config-source/cdb/postgresutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type TokenRegistry struct {
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func NewTokenRegistry(log zerolog.Logger, pool *pgxpool.Pool) *TokenRegistry {
	return &TokenRegistry{
		pool: pool,
		log:  log,
	}
}

type APIToken struct {
	UserID    string    `db:"user_id" json:"-"`
	CreatedAt time.Time `db:"created_at"`
	Token     string    `db:"token"`
}

func (tr *TokenRegistry) IssueAPIToken(ctx context.Context, signingKey []byte, user User) (APIToken, error) {
	tokenString, err := GenerateAPIToken(signingKey, user)
	if err != nil {
		return APIToken{}, err
	}

	return postgresutils.GetOne[APIToken](
		tr.pool,
		ctx,
		"INSERT INTO api_tokens (user_id, token) VALUES ($1, $2) RETURNING *",
		strconv.Itoa(int(user.ID)),
		tokenString,
	)
}

func (tr *TokenRegistry) ListAPITokens(ctx context.Context, user User) ([]APIToken, error) {
	return postgresutils.GetAll[APIToken](
		tr.pool,
		ctx,
		"SELECT created_at,token FROM api_tokens WHERE user_id = $1",
		strconv.Itoa(int(user.ID)),
	)
}

func (tr *TokenRegistry) Revoke(ctx context.Context, token string) error {
	_, err := tr.pool.Query(ctx, "INSERT INTO revoked_tokens (token) VALUES ($1)", token)
	// Means the token is already revoked so there isn't really any error worth
	// reporting.
	if postgresutils.IsUniqueConstraintErr(err) {
		return nil
	}

	return err
}

func (tr *TokenRegistry) IsRevoked(ctx context.Context, token string) (bool, error) {
	var count int
	err := tr.pool.QueryRow(
		ctx,
		"SELECT count(*) FROM revoked_tokens WHERE token = $1",
		token,
	).Scan(&count)
	return count > 0, err
}
