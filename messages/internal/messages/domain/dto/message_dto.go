package dto

import (
	"messages/internal/messages/domain"

	"github.com/google/uuid"
)

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

type MessageResponse struct {
	ID        uuid.UUID `json:"id"`
	ChatID    uuid.UUID `json:"chat_id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Content   string    `json:"content"`
	IsEdited  *bool     `json:"is_edited,omitempty"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at,omitempty"`
}

func (req CreateMessageRequest) ToDomain() domain.CreateMessage {
	return domain.CreateMessage{
		ChatID:   req.ChatID,
		SenderID: req.SenderID,
		Content:  req.Content,
	}
}

func ToMessageResponse(msg domain.Message) MessageResponse {
	return MessageResponse{
		ID:        msg.ID,
		ChatID:    msg.ChatID,
		SenderID:  msg.SenderID,
		Content:   msg.Content,
		IsEdited:  msg.IsEdited,
		CreatedAt: msg.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: msg.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}
}
