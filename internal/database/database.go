package internal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(ctx context.Context, db_url string) (*pgxpool.Pool, error) {
	if db_url == "" {
		return nil, fmt.Errorf("%v file not found or failed to read", db_url)
	}

	pool, err := pgxpool.New(ctx, db_url)

	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %v", err)
	}

	err = pool.Ping(ctx)

	if err != nil {
		return nil, fmt.Errorf("unable to ping, connection failed with DB due to: %v", err)
	}

	fmt.Println("Connection is successful!")

	return pool, nil
}
