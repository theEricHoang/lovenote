package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	config "github.com/theEricHoang/lovenote/backend/internal"
)

var DB *pgxpool.Pool

func InitDB() {
	dsn := config.LoadConfig().DatabaseURL
	ctx := context.Background()

	var err error
	DB, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
