package store

import (
	"messages/internal/messages/store/dbqueries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	*dbqueries.Queries
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{
		Queries: dbqueries.New(db),
		db:      db,
	}
}
