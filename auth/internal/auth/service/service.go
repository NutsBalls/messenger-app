package service

import (
	"context"
	"messenger-app/internal/auth/domain"
	"messenger-app/internal/auth/store"
	"messenger-app/pkg/hasher"
	"messenger-app/pkg/tokens"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTTL  = time.Minute * 30
	RefreshTTL = time.Hour * 24
)

type AuthService struct {
	store  *store.AuthRepository
	secret []byte
}

func NewAuthService(store *store.AuthRepository, secret string) *AuthService {
	return &AuthService{
		store:  store,
		secret: []byte(secret),
	}
}

func (s *AuthService) SignUp(ctx context.Context, user domain.CreateUserParams) (domain.User, error) {
	out, err := s.store.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, nil
	}

	return out, nil
}

func (s *AuthService) Login(ctx context.Context, user1 domain.LogInParams) (domain.LoginData, error) {
	user, err := s.store.GetUserByEmail(ctx, user1.Email)
	if err != nil {
		return domain.LoginData{}, errors.Wrap(err, ErrUserNotFound)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(user1.PasswordHash)); err != nil {
		return domain.LoginData{}, errors.Wrap(err, ErrInvalidToken)
	}

	tokensPair, err := s.GenerateTokens(ctx, user.ID.String())
	if err != nil {
		return domain.LoginData{}, errors.Wrap(err, ErrGenerateTokens)
	}

	cryptedRefreshToken, err := hasher.HashToken(tokensPair.Refresh)
	if err != nil {
		return domain.LoginData{}, errors.Wrap(err, ErrHashToken)
	}

	s.store.LogIn(ctx, user1)

	if err := s.store.UpdateRefreshToken(ctx, user.ID, &cryptedRefreshToken); err != nil {
		return domain.LoginData{}, errors.Wrap(err, ErrUpdateToken)
	}

	loginData := domain.LoginData{
		Email:        user.Email,
		AccessToken:  tokensPair.Access,
		RefreshToken: tokensPair.Refresh,
	}

	return loginData, nil

}

func (s *AuthService) ParseToken(ctx context.Context, token string, expectedType string) (uuid.UUID, error) {
	tokenInfo, err := tokens.Verify(token, s.secret)

	if err != nil {
		return uuid.Nil, errors.Wrap(err, ErrInvalidToken)
	}

	if tokenInfo.Type != expectedType {
		return uuid.Nil, errors.New(ErrTypeToken)
	}

	userID, err := uuid.Parse(tokenInfo.ID)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, ErrParseToken)
	}

	return userID, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, uuid uuid.UUID) (domain.User, error) {
	user, err := s.store.GetUserByID(ctx, uuid)
	if err != nil {
		return domain.User{}, errors.Wrap(err, ErrUserNotFound)
	}

	return user, nil
}

func (s *AuthService) GenerateTokens(ctx context.Context, userID string) (*tokens.TokensPair, error) {
	tokens, err := tokens.GenerateTokens(ctx, userID, AccessTTL, RefreshTTL, string(s.secret))
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *AuthService) UpdateRefreshToken(ctx context.Context, uuid uuid.UUID, refreshToken *string) error {
	if err := s.store.UpdateRefreshToken(ctx, uuid, refreshToken); err != nil {
		return err
	}

	return nil
}
