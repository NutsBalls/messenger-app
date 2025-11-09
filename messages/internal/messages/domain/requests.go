package domain

import "github.com/google/uuid"

type CreateMessageRequest struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
	Content  string
}

type GetMessagesRequest struct {
	ChatID uuid.UUID
}

type EditMessageRequest struct {
	MessageID  uuid.UUID
	NewContent string
}

type DeleteMessageRequest struct {
	MessageID uuid.UUID
}
