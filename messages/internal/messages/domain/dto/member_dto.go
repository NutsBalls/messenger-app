package dto

import (
	"messages/internal/messages/domain"

	"github.com/google/uuid"
)

type UserChatRequest struct {
	ChatID uuid.UUID `json:"chat_id"`
	UserID uuid.UUID `json:"user_id"`
}

func (r UserChatRequest) ToDomain() domain.UserChat {
	return domain.UserChat{
		ChatID: r.ChatID,
		UserID: r.UserID,
	}
}
