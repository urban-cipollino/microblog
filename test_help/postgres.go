package test_help

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func SetupPostgres(t testing.TB) *pgxpool.Pool {
	t.Helper()

	host := env("DB_HOST", "localhost")
	port := env("DB_PORT", "5554")
	user := env("DB_USER", "admin")
	password := env("DB_PASSWORD", "admin")
	dbName := env("DB_NAME", "micro_blog") + "_test"
	sslmode := env("DB_SSLMODE", "disable")

	adminURL := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=%s", user, password, host, port, sslmode)
	testURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbName, sslmode)

	// Create test database if it doesn't exist
	adminDB, err := sql.Open("postgres", adminURL)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	var exists bool
	err = adminDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		adminDB.Close() //nolint:errcheck,gosec
		t.Fatalf("failed to check database existence: %v", err)
	}

	if !exists {
		if _, err := adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
			adminDB.Close() //nolint:errcheck,gosec
			t.Fatalf("failed to create test database: %v", err)
		}
	}
	adminDB.Close() //nolint:errcheck,gosec

	// Run migrations
	migrationDB, err := sql.Open("postgres", testURL)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	_, thisFile, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(thisFile), "..", "migrations")

	if err := goose.Up(migrationDB, migrationsDir); err != nil {
		migrationDB.Close() //nolint:errcheck,gosec
		t.Fatalf("failed to run migrations: %v", err)
	}
	migrationDB.Close() //nolint:errcheck,gosec

	// Create pool
	pool, err := pgxpool.New(context.Background(), testURL)
	if err != nil {
		t.Fatalf("failed to create pgxpool: %v", err)
	}

	// Truncate all tables before test
	truncateTables(t, pool)

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func truncateTables(t testing.TB, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(context.Background(), `
        TRUNCATE post_likes, posts, users RESTART IDENTITY CASCADE
    `)
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}
