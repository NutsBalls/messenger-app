package store

import (
	"context"
	"messenger-app/internal/auth/domain"

	"github.com/google/uuid"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	UpdateRefreshToken(ctx context.Context, id uuid.UUID, refreshToken *string) error
}
