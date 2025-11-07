package store

import (
	"context"
	"messages/internal/messages/domain"
	"messages/internal/messages/store/dbqueries"

	"github.com/google/uuid"
)

func (r *MessageRepository) CreateMessage(ctx context.Context, req domain.CreateMessageRequest) (domain.Message, error) {
	params := dbqueries.CreateMessageParams{
		ChatID:   req.ChatID,
		SenderID: req.SenderID,
		Content:  req.Content,
	}
	dbMsg, err := r.Queries.CreateMessage(ctx, params)
	if err != nil {
		return domain.Message{}, err
	}

	return domain.Message{
		ID:        dbMsg.ID,
		ChatID:    dbMsg.ChatID,
		SenderID:  dbMsg.SenderID,
		Content:   dbMsg.Content,
		IsEdited:  dbMsg.IsEdited,
		CreatedAt: dbMsg.CreatedAt,
		UpdatedAt: dbMsg.UpdatedAt,
	}, nil

}

func (r *MessageRepository) GetMessages(ctx context.Context, chatID uuid.UUID) ([]dbqueries.Message, error) {
	return r.Queries.GetMessages(ctx, chatID)
}

func (r *MessageRepository) EditMessage(ctx context.Context, params dbqueries.EditMessageParams) error {
	return r.Queries.EditMessage(ctx, params)
}

func (r *MessageRepository) DeleteMessage(ctx context.Context, id uuid.UUID) error {
	return r.Queries.DeleteMessage(ctx, id)
}

func (r *MessageRepository) DeleteMessages(ctx context.Context, chatID uuid.UUID) error {
	return r.Queries.DeleteMessages(ctx, chatID)
}

func (r *MessageRepository) MessageExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.Queries.MessageExists(ctx, id)
}
