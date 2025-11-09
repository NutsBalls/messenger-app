package domain

import "github.com/google/uuid"

type CreateMessageRequest struct {
	ChatID   uuid.UUID `json:"chat_id"`
	SenderID uuid.UUID `json:"sender_id"`
	Content  string    `json:"content"`
}

type GetMessagesRequest struct {
	ChatID uuid.UUID `json:"chat_id"`
}

type EditMessageRequest struct {
	MessageID  uuid.UUID `json:"message_id"`
	NewContent string    `json:"new_content"`
}

type DeleteMessageRequest struct {
	MessageID uuid.UUID `json:"message_id"`
}

type CreateChatRequest struct {
	IsGroup bool `json:"is_group"`
}

type DeleteChatRequest struct {
	ChatID uuid.UUID `json:"chat_id"`
}

type CreateGroupChat struct {
	Name string `json:"name"`
}

type UserChatRequest struct {
	ChatID uuid.UUID `json:"chat_id"`
	UserID uuid.UUID `json:"user_id"`
}
