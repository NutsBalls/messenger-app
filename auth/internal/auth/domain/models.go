package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Email               string    `json:"email"`
	IsAdmin             bool      `json:"is_admin"`
	Password            string    `json:"password"`
	CryptedRefreshToken string
}

type CreateUserParams struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password_hash"`
}

type LogInParams struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type LoginData struct {
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
