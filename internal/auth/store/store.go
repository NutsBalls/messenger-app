package store

import (
	"context"
	"messenger-app/internal/auth/store/generated"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
	return r.Queries.GetUser(ctx, pgtype.UUID{
		Bytes: uuid,
		Valid: true,
	})
}

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, user_id uuid.UUID, refresh_token string, expires_at time.Time) (generated.RefreshToken, error) {
	return r.Queries.CreateRefreshToken(ctx, generated.CreateRefreshTokenParams{
		UserID: pgtype.UUID{
			Bytes: user_id,
			Valid: true,
		},
		RefreshToken: refresh_token,
		ExpiresAt:    pgtype.Timestamp{Time: expires_at, Valid: true},
	})
}

func (r *AuthRepository) UpdateRefreshToken(ctx context.Context, user_id uuid.UUID, refresh_token string, expires_at time.Time) (generated.RefreshToken, error) {
	return r.Queries.UpdateRefreshToken(ctx, generated.UpdateRefreshTokenParams{
		UserID: pgtype.UUID{
			Bytes: user_id,
			Valid: true,
		},
		RefreshToken: refresh_token,
		ExpiresAt:    pgtype.Timestamp{Time: expires_at, Valid: true},
	})
}

func (r *AuthRepository) CheckRefreshToken(ctx context.Context, refresh_token string) (bool, error) {
	return r.Queries.CheckRefreshToken(ctx, refresh_token)
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, user_id uuid.UUID) error {
	return r.Queries.DeleteRefreshToken(ctx, pgtype.UUID{
		Bytes: user_id,
		Valid: true,
	})
}

func (r *AuthRepository) LogoutRefreshToken(ctx context.Context, refresh_token string) (generated.RefreshToken, error) {
	return r.Queries.LogoutRefreshToken(ctx, refresh_token)
}
