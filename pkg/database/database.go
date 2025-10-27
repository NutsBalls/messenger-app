package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
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

	return pool
}
