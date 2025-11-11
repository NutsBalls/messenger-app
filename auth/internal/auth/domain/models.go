package domain

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID                  uuid.UUID          `json:"id"`
	Name                string             `json:"name"`
	Email               string             `json:"email"`
	PasswordHash        string             `json:"password_hash"`
	CreatedAt           pgtype.Timestamptz `json:"created_at"`
	UpdatedAt           pgtype.Timestamptz `json:"updated_at"`
	CryptedRefreshToken *string            `json:"crypted_refresh_token"`
}

type CreateUserParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
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
