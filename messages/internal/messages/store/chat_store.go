package store

import (
	"context"
	"messages/internal/messages/store/dbqueries"

	"github.com/google/uuid"
)

func (r *MessageRepository) CreateChat(ctx context.Context, isGroup bool) (dbqueries.Chat, error) {
	return r.Queries.CreateChat(ctx, &isGroup)
}

func (r *MessageRepository) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	return r.Queries.DeleteChat(ctx, chatID)
}

func (r *MessageRepository) CreateGroupChat(ctx context.Context, name *string) (dbqueries.Chat, error) {
	return r.Queries.CreateGroupChat(ctx, name)
}

func (r *MessageRepository) ChatExists(ctx context.Context, chatID uuid.UUID) (bool, error) {
	return r.Queries.ChatExists(ctx, chatID)
}

func (r *MessageRepository) IsUserInChat(ctx context.Context, params dbqueries.IsUserInChatParams) (bool, error) {
	return r.Queries.IsUserInChat(ctx, params)
}
