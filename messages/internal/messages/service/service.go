package service

import (
	"context"
	"messages/internal/messages/domain"
	"messages/internal/messages/store"
	"messages/internal/messages/store/dbqueries"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type MessageService struct {
	store *store.MessageRepository
}

func NewMessageService(store *store.MessageRepository) *MessageService {
	return &MessageService{
		store: store,
	}
}

// CreateMessage ..
func (s *MessageService) CreateMessage(ctx context.Context, chatID uuid.UUID, senderID uuid.UUID, content string) (domain.Message, error) {
	exist, err := s.store.ChatExists(ctx, chatID)
	if err != nil {
		return domain.Message{}, errors.Wrap(err, "internal store")
	}
	if !exist {
		return domain.Message{}, errors.Wrap(err, "chat doesn't exists")
	}

	inChat, err := s.store.IsUserInChat(ctx, dbqueries.IsUserInChatParams{
		ChatID: chatID,
		UserID: senderID,
	})
	if err != nil {
		return domain.Message{}, errors.Wrap(err, "internal store")
	}
	if !inChat {
		return domain.Message{}, errors.Wrap(err, "user is not in the chat")
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

// GetMessages
func (s *MessageService) GetMessages(ctx context.Context, chatID uuid.UUID) ([]dbqueries.Message, error) {
	exist, err := s.store.ChatExists(ctx, chatID)
	if err != nil {
		return []dbqueries.Message{}, errors.Wrap(err, "internal store")
	}
	if !exist {
		return []dbqueries.Message{}, errors.Wrap(err, "chat doesn't exists")
	}

	msgs, err := s.store.GetMessages(ctx, chatID)
	if err != nil {
		return []dbqueries.Message{}, errors.Wrap(err, "internal store")
	}

	return msgs, nil
}

// EditMessage
func (s *MessageService) EditMessage(ctx context.Context, messageID uuid.UUID, newContent string) error {
	exist, err := s.store.MessageExists(ctx, messageID)
	if err != nil {
		return errors.Wrap(err, "internal store")
	}
	if !exist {
		return errors.Wrap(err, "message doesn't exists")
	}

	err = s.store.EditMessage(ctx, dbqueries.EditMessageParams{
		ID:      messageID,
		Content: newContent,
	})
	if err != nil {
		return errors.Wrap(err, "internal store")
	}

	return nil
}

// DeleteMessage
func (s *MessageService) DeleteMessage(ctx context.Context, messageID uuid.UUID) error {
	exist, err := s.store.MessageExists(ctx, messageID)
	if err != nil {
		return errors.Wrap(err, "internal store")
	}
	if !exist {
		return errors.Wrap(err, "message doesn't exists")
	}

	err = s.store.DeleteMessage(ctx, messageID)
	if err != nil {
		return errors.Wrap(err, "internal store")
	}

	return nil
}
