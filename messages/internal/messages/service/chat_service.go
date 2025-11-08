package service

import (
	"context"
	"messages/internal/messages/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreateChat create chat with user to user
func (s *MessagesService) CreateChat(ctx context.Context, isGroup bool) (domain.Chat, error) {
	chat, err := s.store.CreateChat(ctx, isGroup)
	if err != nil {
		return domain.Chat{}, errors.Wrap(err, "internal store")
	}

	return chat, nil
}

// DeleteChat delete chat and messages in chat
func (s *MessagesService) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	exists, err := s.store.ChatExists(ctx, chatID)
	if err != nil {
		return errors.Wrap(err, "internal store")
	}
	if !exists {
		return domain.ErrChatNotFound
	}

	if err := s.store.DeleteChat(ctx, chatID); err != nil {
		return errors.Wrap(err, "internal store")
	}

	if err := s.store.DeleteMessages(ctx, chatID); err != nil {
		return errors.Wrap(err, "internal store")
	}

	return nil
}

// CreateGroupChat create chat with user to users
func (s *MessagesService) CreateGroupChat(ctx context.Context, name *string) (domain.Chat, error) {
	groupChat, err := s.store.CreateGroupChat(ctx, name)
	if err != nil {
		return domain.Chat{}, errors.Wrap(err, "internal store")
	}

	return groupChat, nil
}
