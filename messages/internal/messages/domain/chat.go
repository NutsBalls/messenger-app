package domain

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Chat struct {
	ID        uuid.UUID
	Name      string
	IsGroup   bool
	CreatedAt pgtype.Timestamptz
}
