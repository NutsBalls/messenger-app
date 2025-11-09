package domain

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateMessage struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
	Content  string
}

type Message struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	Content   string
	IsEdited  *bool
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}
