package service

import (
	"context"
	"messages/internal/messages/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// AddUserToChat add member to a created chat
func (s *MessagesService) AddUserToChat(ctx context.Context, params domain.UserChat) error {
	chatExists, err := s.store.ChatExists(ctx, params.ChatID)
	if err != nil {
		return errors.Wrap(err, "internal store")
	}
	if !chatExists {
		return domain.ErrChatNotFound
	}

	if err := s.store.AddUserToChat(ctx, params); err != nil {
		return errors.Wrap(err, "internal store")
	}

	return nil
}

// RemoveUserFromChat removes a member from a created chat
func (s *MessagesService) RemoveUserFromChat(ctx context.Context, params domain.UserChat) error {
	chatExists, err := s.store.ChatExists(ctx, params.ChatID)
	if err != nil {
		return errors.Wrap(err, "internal store")
	}
	if !chatExists {
		return domain.ErrChatNotFound
	}

	inChat, err := s.store.IsUserInChat(ctx, params.ChatID, params.UserID)
	if err != nil {
		return errors.Wrap(err, "internal store")
	}
	if !inChat {
		return domain.ErrUserNotInChat
	}

	if err := s.store.RemoveUserFromChat(ctx, params); err != nil {
		return errors.Wrap(err, "internal store")
	}

	return nil
}

// GetChatMembers shows all chat members
func (s *MessagesService) GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	chatExists, err := s.store.ChatExists(ctx, chatID)
	if err != nil {
		return []uuid.UUID{}, errors.Wrap(err, "internal store")
	}
	if !chatExists {
		return []uuid.UUID{}, domain.ErrChatNotFound
	}

	members, err := s.store.GetChatMembers(ctx, chatID)
	if err != nil {
		return []uuid.UUID{}, errors.Wrap(err, "internal store")
	}

	return members, nil
}

// GetUserChats shows all chats that the user has
func (s *MessagesService) GetUserChats(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error) {
	chats, err := s.store.GetUserChats(ctx, userID)
	if err != nil {
		return []domain.Chat{}, errors.Wrap(err, "internal store")
	}

	return chats, nil
}
