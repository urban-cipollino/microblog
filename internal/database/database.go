package database

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	"micro-blog/internal/config"
)

var Set = wire.NewSet(
	NewPool,
	NewQueryBuilder,
)

func NewPool(databaseURL config.DatabaseURL) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, string(databaseURL))
	if err != nil {
		panic(err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		panic(err)
	}

	return pool
}

func NewQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
