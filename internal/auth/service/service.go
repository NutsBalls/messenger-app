package service

import (
	"context"
	"messenger-app/internal/auth/store"
	"messenger-app/internal/auth/store/generated"
	"messenger-app/pkg/hasher"
	"messenger-app/pkg/tokens"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	ErrWithLogin              = "smth wrong in login"
	ErrWithUpdateRefreshToken = "smth wrong in create refresh token"
	AccessTTL                 = time.Minute * 30
	RefreshTTL                = time.Hour * 24
)

type AuthService struct {
	store  *store.AuthRepository
	secret []byte
}

type LoginData struct {
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthService(store *store.AuthRepository, secret string) *AuthService {
	return &AuthService{
		store:  store,
		secret: []byte(secret),
	}
}

func (s *AuthService) SignUp(ctx context.Context, user generated.CreateUserParams) (generated.User, error) {
	return s.store.CreateUser(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, user1 generated.LogInParams) (*LoginData, error) {
	user, err := s.store.GetUserByEmail(ctx, user1.Email)
	if err != nil {
		return nil, errors.Wrap(err, ErrWithLogin)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(user1.PasswordHash)); err != nil {
		return nil, errors.Wrap(err, ErrWithLogin)
	}

	tokensPair, err := s.GenerateTokens(ctx, user.ID.String())
	if err != nil {
		return nil, errors.Wrap(err, ErrWithLogin)
	}

	cryptedRefreshToken, err := hasher.HashToken(tokensPair.Refresh)
	if err != nil {
		return nil, errors.Wrap(err, ErrWithLogin)
	}

	s.store.LogIn(ctx, user1)

	if err := s.store.UpdateRefreshToken(ctx, user.ID, &cryptedRefreshToken); err != nil {
		return nil, errors.Wrap(err, "failed to delete old refresh token")
	}

	return &LoginData{
		Email:        user.Email,
		AccessToken:  tokensPair.Access,
		RefreshToken: tokensPair.Refresh,
	}, nil

}

func (s *AuthService) ParseToken(ctx context.Context, token string, expectedType string) (uuid.UUID, error) {
	tokenInfo, err := tokens.Verify(token, s.secret)

	if err != nil {
		return uuid.Nil, errors.New("invalid token after verify")
	}

	if tokenInfo.Type != expectedType {
		return uuid.Nil, errors.New("invalid token type")
	}

	userID, err := uuid.Parse(tokenInfo.ID)
	if err != nil {
		return uuid.Nil, errors.New("invalid token after parse")
	}

	return userID, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, uuid uuid.UUID) (generated.User, error) {
	user, err := s.store.GetUserByID(ctx, uuid)
	if err != nil {
		return generated.User{}, errors.New("invalid token in get user by ID")
	}

	return user, nil
}

func (s *AuthService) GenerateTokens(ctx context.Context, userID string) (*tokens.TokensPair, error) {
	return tokens.GenerateTokens(ctx, userID, AccessTTL, RefreshTTL, string(s.secret))
}

func (s *AuthService) UpdateRefreshToken(ctx context.Context, uuid uuid.UUID, refreshToken *string) error {
	return s.store.UpdateRefreshToken(ctx, uuid, refreshToken)
}
