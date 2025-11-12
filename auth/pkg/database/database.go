package database

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func NewConn(databaseURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	SetupMigrations(pool)
	return pool
}

func SetupMigrations(pool *pgxpool.Pool) {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	wd, _ := os.Getwd()
	migrationsDir := filepath.Join(wd, "migrations")

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}
