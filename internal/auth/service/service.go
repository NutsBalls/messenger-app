package service

import (
	"context"
	"messenger-app/internal/auth/store"
	"messenger-app/internal/auth/store/generated"
	"messenger-app/pkg/tokens"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	ErrWithLogin              = "smth wrong in login"
	ErrWithUpdateRefreshToken = "smth wrong in create refresh token"
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

func (s *AuthService) SignUp(ctx context.Context, user generated.CreateUserParams) (generated.User, error) {
	return s.store.CreateUser(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, user1 generated.LogInParams) (string, string, generated.User, error) {
	user, err := s.store.GetUserByEmail(ctx, user1.Email)
	if err != nil {
		return "", "", generated.User{}, errors.Wrap(err, ErrWithLogin)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(user1.PasswordHash)); err != nil {
		return "", "", generated.User{}, errors.Wrap(err, ErrWithLogin)
	}

	accessToken, err := s.GenerateAccessToken(ctx, user.ID.String(), time.Duration(time.Minute*30))
	if err != nil {
		return "", "", generated.User{}, errors.Wrap(err, ErrWithLogin)
	}

	refreshToken, err := s.GenerateRefreshToken(ctx, user.ID.String(), time.Duration(time.Hour*24))
	if err != nil {
		return "", "", generated.User{}, errors.Wrap(err, ErrWithLogin)
	}

	s.store.LogIn(ctx, user1)
	expiresAt := time.Now().Add(24 * time.Hour)

	if err := s.store.DeleteRefreshToken(ctx, user.ID.Bytes); err != nil {
		return "", "", generated.User{}, errors.Wrap(err, "failed to delete old refresh token")
	}

	_, err = s.store.CreateRefreshToken(ctx, user.ID.Bytes, refreshToken, expiresAt)
	if err != nil {
		return "", "", generated.User{}, errors.Wrap(err, "failed to create new refresh token")
	}

	return accessToken, refreshToken, user, nil

}

func (s *AuthService) GenerateAccessToken(ctx context.Context, userID string, ttl time.Duration) (string, error) {
	claims := tokens.TokenClaims{
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "messenger-app",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *AuthService) GenerateRefreshToken(ctx context.Context, userID string, ttl time.Duration) (string, error) {
	claims := tokens.TokenClaims{
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "messenger-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)

}

func (s *AuthService) CheckRefreshToken(ctx context.Context, refresh_token string) (bool, error) {
	check, err := s.store.CheckRefreshToken(ctx, refresh_token)
	if !check || err != nil {
		return false, errors.Wrap(err, ErrWithUpdateRefreshToken)
	}

	return check, nil
}

func (s *AuthService) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, expiresAt time.Time) error {
	_, err := s.store.UpdateRefreshToken(ctx, userID, refreshToken, expiresAt)
	if err != nil {
		return errors.Wrap(err, ErrWithUpdateRefreshToken)
	}
	return nil
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
