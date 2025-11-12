package service

import (
	"context"
	"errors"
	"messenger-app/internal/auth/domain"
	"messenger-app/internal/auth/store"
	"messenger-app/pkg/hasher"
	"messenger-app/pkg/tokens"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTTL  = 30 * time.Minute
	RefreshTTL = 24 * time.Hour
)

type AuthService struct {
	repo   store.AuthRepository
	secret []byte
}

func NewAuthService(repo store.AuthRepository, secret string) *AuthService {
	return &AuthService{
		repo:   repo,
		secret: []byte(secret),
	}
}

func (s *AuthService) SignUp(ctx context.Context, user domain.CreateUserParams) (domain.User, error) {
	hashedPassword, err := hasher.Hash(user.Password)
	if err != nil {
		return domain.User{}, err
	}

	user.Password = hashedPassword
	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, creds domain.LogInParams) (domain.LoginData, error) {
	user, err := s.repo.GetUserByEmail(ctx, creds.Email)
	if err != nil {
		return domain.LoginData{}, errors.New(ErrUserNotFound)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.PasswordHash)); err != nil {
		return domain.LoginData{}, errors.New(ErrInvalidCredentials)
	}

	tokensPair, err := s.GenerateTokens(ctx, user.ID.String())
	if err != nil {
		return domain.LoginData{}, errors.New(ErrGenerateTokens)
	}

	hashedRefresh, err := hasher.HashToken(tokensPair.Refresh)
	if err != nil {
		return domain.LoginData{}, errors.New(ErrHashToken)
	}

	if err := s.repo.UpdateRefreshToken(ctx, user.ID, &hashedRefresh); err != nil {
		return domain.LoginData{}, errors.New(ErrUpdateToken)
	}

	return domain.LoginData{
		Email:        user.Email,
		AccessToken:  tokensPair.Access,
		RefreshToken: tokensPair.Refresh,
	}, nil
}

func (s *AuthService) ParseToken(ctx context.Context, token string, expectedType string) (uuid.UUID, error) {
	info, err := tokens.Verify(token, s.secret)
	if err != nil {
		return uuid.Nil, errors.New(ErrInvalidToken)
	}

	if info.Type != expectedType {
		return uuid.Nil, errors.New(ErrTypeToken)
	}

	id, err := uuid.Parse(info.ID)
	if err != nil {
		return uuid.Nil, errors.New(ErrParseToken)
	}

	return id, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *AuthService) GenerateTokens(ctx context.Context, userID string) (*tokens.TokensPair, error) {
	return tokens.GenerateTokens(ctx, userID, AccessTTL, RefreshTTL, string(s.secret))
}

func (s *AuthService) UpdateRefreshToken(ctx context.Context, id uuid.UUID, refreshToken *string) error {
	return s.repo.UpdateRefreshToken(ctx, id, refreshToken)
}
