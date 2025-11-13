package store

import (
	"context"
	"messages/internal/messages/domain"

	"github.com/google/uuid"
)

type Repository interface {
	CreateMessage(ctx context.Context, req domain.CreateMessage) (domain.Message, error)
	GetMessages(ctx context.Context, chatID uuid.UUID) ([]domain.Message, error)
	EditMessage(ctx context.Context, msgID uuid.UUID, newContent string) error
	DeleteMessage(ctx context.Context, msgID uuid.UUID) error
	MessageExists(ctx context.Context, msgID uuid.UUID) (bool, error)

	CreateChat(ctx context.Context, isGroup bool) (domain.Chat, error)
	DeleteChat(ctx context.Context, chatID uuid.UUID) error
	CreateGroupChat(ctx context.Context, name *string) (domain.Chat, error)
	ChatExists(ctx context.Context, chatID uuid.UUID) (bool, error)
	DeleteMessages(ctx context.Context, chatID uuid.UUID) error
	IsUserInChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (bool, error)

	AddUserToChat(ctx context.Context, params domain.UserChat) error
	RemoveUserFromChat(ctx context.Context, params domain.UserChat) error
	GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error)
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error)
}
