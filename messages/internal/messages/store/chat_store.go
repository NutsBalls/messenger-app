package store

import (
	"context"
	"messages/internal/messages/domain"
	"messages/internal/messages/store/dbqueries"

	"github.com/google/uuid"
)

func (r *MessagesRepository) CreateChat(ctx context.Context, isGroup bool) (domain.Chat, error) {
	dbChat, err := r.Queries.CreateChat(ctx, &isGroup)
	if err != nil {
		return domain.Chat{}, err
	}

	var name string
	if dbChat.Name != nil {
		name = *dbChat.Name
	}
	if dbChat.IsGroup != nil {
		isGroup = *dbChat.IsGroup
	}

	chat := domain.Chat{
		ID:        dbChat.ID,
		Name:      name,
		IsGroup:   isGroup,
		CreatedAt: dbChat.CreatedAt,
	}

	return chat, nil
}

func (r *MessagesRepository) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	if err := r.Queries.DeleteChat(ctx, chatID); err != nil {
		return err
	}
	return nil
}

func (r *MessagesRepository) CreateGroupChat(ctx context.Context, name *string) (domain.Chat, error) {
	params, err := r.Queries.CreateGroupChat(ctx, name)
	if err != nil {
		return domain.Chat{}, err
	}

	groupChat := domain.Chat{
		ID:        params.ID,
		Name:      *params.Name,
		IsGroup:   *params.IsGroup,
		CreatedAt: params.CreatedAt,
	}

	return groupChat, nil
}

func (r *MessagesRepository) ChatExists(ctx context.Context, chatID uuid.UUID) (bool, error) {
	exists, err := r.Queries.ChatExists(ctx, chatID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *MessagesRepository) IsUserInChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error) {
	params := dbqueries.IsUserInChatParams{
		ChatID: chatID,
		UserID: userID,
	}

	exists, err := r.Queries.IsUserInChat(ctx, params)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *MessagesRepository) DeleteMessages(ctx context.Context, chatID uuid.UUID) error {
	if err := r.Queries.DeleteMessages(ctx, chatID); err != nil {
		return err
	}

	return nil
}
