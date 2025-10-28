package tokens

import (
	"errors"

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
