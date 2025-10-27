package service

import (
	"context"
	"messenger-app/internal/auth/store"
	"messenger-app/internal/auth/store/generated"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func NewAuthService(store *store.AuthRepository) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (s *AuthService) SignUp(ctx context.Context, user generated.CreateUserParams) (generated.User, error) {
	return s.store.CreateUser(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, user1 generated.LogInParams) (string, generated.User, error) {
	user, err := s.store.GetUserByEmail(context.Background(), user1.Email)
	if err != nil {
		return "", generated.User{}, errors.Wrap(err, ErrWithLogin)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(user1.PasswordHash)); err != nil {
		return "", generated.User{}, errors.Wrap(err, ErrWithLogin)
	}

	token, err := GenerateAccessToken(user.ID.String(), s.secret, time.Duration(time.Minute)*30)
	if err != nil {
		return "", generated.User{}, errors.Wrap(err, ErrWithLogin)
	}

	s.store.LogIn(ctx, user1)

	return token, user, nil

}

func GenerateAccessToken(userID string, secret []byte, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
		"iss": "messenger-app",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
