package store

import (
	"context"
	"messages/internal/messages/store/dbqueries"

	"github.com/google/uuid"
)

func (r *MessageRepository) AddUserToChat(ctx context.Context, params dbqueries.AddUserToChatParams) error {
	return r.Queries.AddUserToChat(ctx, params)
}

func (r *MessageRepository) RemoveUserFromChat(ctx context.Context, params dbqueries.RemoveUserFromChatParams) error {
	return r.Queries.RemoveUserFromChat(ctx, params)
}

func (r *MessageRepository) GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	return r.Queries.GetChatMembers(ctx, chatID)
}

func (r *MessageRepository) GetUserChats(ctx context.Context, userID uuid.UUID) ([]dbqueries.Chat, error) {
	return r.Queries.GetUserChats(ctx, userID)
}
