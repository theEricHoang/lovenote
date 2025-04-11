package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	config "github.com/theEricHoang/lovenote/backend/internal"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase() (*Database, error) {
	dsn := config.LoadConfig().DatabaseURL
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
