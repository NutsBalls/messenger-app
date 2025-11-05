package tokens

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

type VerificationResponse struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type TokensPair struct {
	Access, Refresh string
}

func Verify(token string, secret []byte) (*VerificationResponse, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*TokenClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token claims")
	}

	return &VerificationResponse{
		Type: claims.Type,
		ID:   claims.Subject,
	}, nil
}

func GenerateTokens(ctx context.Context, userID string, accessTTL, refreshTTL time.Duration, secret string) (*TokensPair, error) {
	accessClaims := TokenClaims{
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "messenger-app",
		},
	}

	refreshClaims := TokenClaims{
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "messenger-app",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	refreshString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &TokensPair{
		Access:  accessString,
		Refresh: refreshString,
	}, nil
}
