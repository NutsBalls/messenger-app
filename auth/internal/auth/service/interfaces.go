package service

import (
	"context"
	"messenger-app/internal/auth/domain"
	"messenger-app/pkg/tokens"

	"github.com/google/uuid"
)

type AuthUsecase interface {
	SignUp(ctx context.Context, params domain.CreateUserParams) (domain.User, error)
	Login(ctx context.Context, params domain.LogInParams) (domain.LoginData, error)
	ParseToken(ctx context.Context, token, kind string) (uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GenerateTokens(ctx context.Context, userID string) (*tokens.TokensPair, error)
	UpdateRefreshToken(ctx context.Context, userID uuid.UUID, token *string) error
}
