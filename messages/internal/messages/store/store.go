package store

import (
	"messages/internal/messages/store/dbqueries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessagesRepository struct {
	*dbqueries.Queries
	db *pgxpool.Pool
}

func NewMessagesRepository(db *pgxpool.Pool) *MessagesRepository {
	return &MessagesRepository{
		Queries: dbqueries.New(db),
		db:      db,
	}
}
