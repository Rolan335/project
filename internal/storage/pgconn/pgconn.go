package pgconn

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConn(connStr string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres.NewConn: %w", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("postgres.NewConn: %w", err)
	}

	return conn, nil
}
