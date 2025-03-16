package pgconn

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func GetConn(connStr string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, errors.Wrap(err, "postgres.NewConn")
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, errors.Wrap(err, "postgres.NewConn")
	}

	return conn, nil
}
