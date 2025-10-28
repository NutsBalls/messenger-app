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

var (
	ErrWithLogin string
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

func (s *AuthService) Login(ctx context.Context, user1 generated.LogInParams) (string, generated.User, error) {
	user, err := s.store.GetUserByEmail(ctx, user1.Email)
	if err != nil {
		return "", generated.User{}, errors.Wrap(err, ErrWithLogin)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(user1.PasswordHash)); err != nil {
		return "", generated.User{}, errors.Wrap(err, ErrWithLogin)
	}

	token, err := s.GenerateAccessToken(user.ID.String(), time.Duration(time.Minute)*30)
	if err != nil {
		return "", generated.User{}, errors.Wrap(err, ErrWithLogin)
	}

	s.store.LogIn(ctx, user1)

	return token, user, nil

}

func (s *AuthService) GenerateAccessToken(userID string, ttl time.Duration) (string, error) {
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

func (s *AuthService) ParseToken(accesToken string, expectedType string) (uuid.UUID, error) {
	tokenInfo, err := tokens.Verify(accesToken, s.secret)

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
