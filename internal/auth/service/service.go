package service

import (
	"context"
	"messenger-app/internal/auth/store"
	"messenger-app/internal/auth/store/generated"
)

type AuthService struct {
	store *store.AuthRepository
}

func NewAuthService(store *store.AuthRepository) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user generated.CreateUserParams) (generated.User, error) {
	return s.store.CreateUser(ctx, user)
}
