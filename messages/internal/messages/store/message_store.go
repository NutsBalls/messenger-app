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

func (r *MessageRepository) GetMessages(ctx context.Context, chatID uuid.UUID) ([]domain.Message, error) {
	dbMsgs, err := r.Queries.GetMessages(ctx, chatID)
	if err != nil {
		return []domain.Message{}, err
	}

	msgs := make([]domain.Message, 0, len(dbMsgs))
	for _, v := range dbMsgs {
		msgs = append(msgs, domain.Message{
			ID:        v.ID,
			ChatID:    v.ChatID,
			SenderID:  v.SenderID,
			Content:   v.Content,
			IsEdited:  v.IsEdited,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return msgs, err
}

func (r *MessageRepository) EditMessage(ctx context.Context, msgID uuid.UUID, newContent string) error {
	params := dbqueries.EditMessageParams{
		ID:      msgID,
		Content: newContent,
	}

	if err := r.Queries.EditMessage(ctx, params); err != nil {
		return err
	}

	return nil
}

func (r *MessageRepository) DeleteMessage(ctx context.Context, msgID uuid.UUID) error {
	if err := r.Queries.DeleteMessage(ctx, msgID); err != nil {
		return err
	}

	return nil
}

func (r *MessageRepository) MessageExists(ctx context.Context, msgID uuid.UUID) (bool, error) {
	exists, err := r.Queries.MessageExists(ctx, msgID)
	if err != nil {
		return false, err
	}
	return exists, nil
}
