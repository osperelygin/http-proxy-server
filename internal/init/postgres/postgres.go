package init_postgres

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(ctx context.Context, env string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, os.Getenv(env))
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
