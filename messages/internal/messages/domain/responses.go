package domain

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateMessageResponse struct {
	ChatID    uuid.UUID          `json:"chat_id"`
	SenderID  uuid.UUID          `json:"sender_id"`
	Content   string             `json:"content"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type ChatResponse struct {
	Name    string `json:"name"`
	IsGroup bool   `json:"is_group"`
}
