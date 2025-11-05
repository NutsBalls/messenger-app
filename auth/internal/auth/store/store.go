package store

import (
	"context"
	"messenger-app/internal/auth/store/generated"

	"github.com/google/uuid"
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

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (generated.User, error) {
	return r.Queries.GetUserByEmail(ctx, email)
}

func (r *AuthRepository) LogIn(ctx context.Context, params generated.LogInParams) (generated.User, error) {
	return r.Queries.LogIn(ctx, params)
}

func (r *AuthRepository) GetUserByID(ctx context.Context, uuid uuid.UUID) (generated.User, error) {
	return r.Queries.GetUser(ctx, uuid)
}

func (r *AuthRepository) UpdateRefreshToken(ctx context.Context, uuid uuid.UUID, refreshToken *string) error {
	return r.Queries.UpdateRefreshToken(ctx, generated.UpdateRefreshTokenParams{
		ID:                  uuid,
		CryptedRefreshToken: refreshToken,
	})
}
