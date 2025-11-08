package service

import (
	"context"
	"messages/internal/messages/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreateMessage creates a message in the created chat
func (s *MessagesService) CreateMessage(ctx context.Context, chatID uuid.UUID, senderID uuid.UUID, content string) (domain.Message, error) {
	exist, err := s.store.ChatExists(ctx, chatID)
	if err != nil {
		return domain.Message{}, errors.Wrap(err, "internal store")
	}
	if !exist {
		return domain.Message{}, domain.ErrChatNotFound
	}

	inChat, err := s.store.IsUserInChat(ctx, chatID, senderID)
	if err != nil {
		return domain.Message{}, errors.Wrap(err, "internal store")
	}
	if !inChat {
		return domain.Message{}, domain.ErrUserNotInChat
	}

	msg, err := s.store.CreateMessage(ctx, domain.CreateMessageRequest{
		ChatID:   chatID,
		SenderID: senderID,
		Content:  content,
	})
	if err != nil {
		return domain.Message{}, errors.Wrap(err, "internal store")
	}

	return msg, nil
}

// GetMessages receives messages from the chat
func (s *MessagesService) GetMessages(ctx context.Context, chatID uuid.UUID) ([]domain.Message, error) {
	exist, err := s.store.ChatExists(ctx, chatID)
	if err != nil {
		return []domain.Message{}, errors.Wrap(err, "internal store")
	}
	if !exist {
		return []domain.Message{}, domain.ErrChatNotFound
	}

	msgs, err := s.store.GetMessages(ctx, chatID)
	if err != nil {
		return []domain.Message{}, errors.Wrap(err, "internal store")
	}

	return msgs, nil
}

// EditMessage changes user message
func (s *MessagesService) EditMessage(ctx context.Context, msgID uuid.UUID, newContent string) error {
	exist, err := s.store.MessageExists(ctx, msgID)
	if err != nil {
		return errors.Wrap(err, "internal store")
	}
	if !exist {
		return domain.ErrMessageNotFound
	}

	err = s.store.EditMessage(ctx, msgID, newContent)
	if err != nil {
		return errors.Wrap(err, "internal store")
	}

	return nil
}

// DeleteMessage delete user message
func (s *MessagesService) DeleteMessage(ctx context.Context, msgID uuid.UUID) error {
	exist, err := s.store.MessageExists(ctx, msgID)
	if err != nil {
		return errors.Wrap(err, "internal store")
	}
	if !exist {
		return domain.ErrMessageNotFound
	}

	err = s.store.DeleteMessage(ctx, msgID)
	if err != nil {
		return errors.Wrap(err, "internal store")
	}

	return nil
}
