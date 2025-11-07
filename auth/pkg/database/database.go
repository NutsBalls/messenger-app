package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func NewConn(databaseURL string) *pgxpool.Pool {
	log.Printf("Connecting to database with URL: %s", databaseURL)
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		fmt.Println("smth wrong")
		log.Fatal(err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	SetupMigrations(pool)

	return pool
}

func SetupMigrations(pool *pgxpool.Pool) {

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	db := stdlib.OpenDBFromPool(pool)
	wd, _ := os.Getwd()
	migrationsDir := filepath.Join(wd, "migrations")

	if err := goose.Up(db, migrationsDir); err != nil {
		panic(err)
	}
	if err := db.Close(); err != nil {
		panic(err)
	}
}
