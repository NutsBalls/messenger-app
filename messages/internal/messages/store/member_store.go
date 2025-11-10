package store

import (
	"context"
	"messages/internal/messages/domain"
	"messages/internal/messages/store/dbqueries"

	"github.com/google/uuid"
)

func (r *MessagesRepository) AddUserToChat(ctx context.Context, params domain.UserChat) error {
	dbParams := dbqueries.AddUserToChatParams{
		UserID: params.UserID,
		ChatID: params.ChatID,
	}
	if err := r.Queries.AddUserToChat(ctx, dbParams); err != nil {
		return err
	}

	return nil
}

func (r *MessagesRepository) RemoveUserFromChat(ctx context.Context, params domain.UserChat) error {
	remove := dbqueries.RemoveUserFromChatParams{
		ChatID: params.ChatID,
		UserID: params.UserID,
	}

	if err := r.Queries.RemoveUserFromChat(ctx, remove); err != nil {
		return err
	}

	return nil
}

func (r *MessagesRepository) GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	users, err := r.Queries.GetChatMembers(ctx, chatID)
	if err != nil {
		return []uuid.UUID{}, err
	}

	return users, nil
}

func (r *MessagesRepository) GetUserChats(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error) {
	dbParams, err := r.Queries.GetUserChats(ctx, userID)
	if err != nil {
		return []domain.Chat{}, err
	}

	chats := make([]domain.Chat, 0, len(dbParams))
	for _, v := range dbParams {
		chats = append(chats, domain.Chat{
			ID:        v.ID,
			Name:      *v.Name,
			IsGroup:   *v.IsGroup,
			CreatedAt: v.CreatedAt,
		})
	}

	return chats, nil
}
