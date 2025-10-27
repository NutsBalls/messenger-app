package store

import (
	"context"
	"messenger-app/internal/auth/store/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	*generated.Queries
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		Queries: generated.New(db),
		db:      db,
	}
}

func (r *AuthRepository) CreateUser(ctx context.Context, params generated.CreateUserParams) (generated.User, error) {
	return r.Queries.CreateUser(ctx, params)
}
