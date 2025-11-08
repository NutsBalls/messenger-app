package domain

import "github.com/google/uuid"

type UserChat struct {
	ChatID uuid.UUID
	UserID uuid.UUID
}
