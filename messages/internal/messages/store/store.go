package store

import (
	"context"
	"messages/internal/messages/store/store"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	*store.Queries
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		Queries: store.New(db),
		db:      db,
	}
}

func (r *AuthRepository) CreateMessage(ctx context.Context, params store.CreateMessageParams) (store.Message, error) {
	return r.Queries.CreateMessage(ctx, params)
}

func (r *AuthRepository) GetMessages(ctx context.Context, chatID uuid.UUID) ([]store.Message, error) {
	return r.Queries.GetMessages(ctx, chatID)
}

func (r *AuthRepository) EditMessage(ctx context.Context, params store.EditMessageParams) error {
	return r.Queries.EditMessage(ctx, params)
}

func (r *AuthRepository) DeleteMessage(ctx context.Context, id uuid.UUID) error {
	return r.Queries.DeleteMessage(ctx, id)
}

func (r *AuthRepository) CreateChat(ctx context.Context, isGroup bool) (store.Chat, error) {
	return r.Queries.CreateChat(ctx, &isGroup)
}

func (r *AuthRepository) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	return r.Queries.DeleteChat(ctx, chatID)
}

func (r *AuthRepository) DeleteMessages(ctx context.Context, chatID uuid.UUID) error {
	return r.Queries.DeleteMessages(ctx, chatID)
}

func (r *AuthRepository) CreateGroupChat(ctx context.Context, name *string) (store.Chat, error) {
	return r.Queries.CreateGroupChat(ctx, name)
}

func (r *AuthRepository) AddUserToChat(ctx context.Context, params store.AddUserToChatParams) error {
	return r.Queries.AddUserToChat(ctx, params)
}

func (r *AuthRepository) RemoveUserFromChat(ctx context.Context, params store.RemoveUserFromChatParams) error {
	return r.Queries.RemoveUserFromChat(ctx, params)
}

func (r *AuthRepository) GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	return r.Queries.GetChatMembers(ctx, chatID)
}

func (r *AuthRepository) GetUserChats(ctx context.Context, userID uuid.UUID) ([]store.Chat, error) {
	return r.Queries.GetUserChats(ctx, userID)
}
