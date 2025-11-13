package dto

import (
	"messages/internal/messages/domain"

	"github.com/google/uuid"
)

type CreateChatRequest struct {
	IsGroup bool `json:"is_group"`
}

type DeleteChatRequest struct {
	ChatID uuid.UUID `json:"chat_id"`
}

type CreateGroupChatRequest struct {
	Name string `json:"name"`
}

type ChatResponse struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	IsGroup bool      `json:"is_group"`
}

func ToChatResponse(chat domain.Chat) ChatResponse {
	return ChatResponse{
		ID:      chat.ID,
		Name:    chat.Name,
		IsGroup: chat.IsGroup,
	}
}
