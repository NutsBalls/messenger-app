package store

import (
	"context"
	"messenger-app/internal/auth/domain"
	"messenger-app/internal/auth/store/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthStore struct {
	*generated.Queries
	db *pgxpool.Pool
}

func NewAuthStore(db *pgxpool.Pool) *AuthStore {
	return &AuthStore{
		Queries: generated.New(db),
		db:      db,
	}
}

func (r *AuthStore) CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
	dbUser, err := r.Queries.CreateUser(ctx, generated.CreateUserParams{
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: params.Password,
	})
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		ID:    dbUser.ID,
		Name:  dbUser.Name,
		Email: dbUser.Email,
	}
	return user, nil
}

func (r *AuthStore) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	dbUser, err := r.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{
		ID:       dbUser.ID,
		Name:     dbUser.Name,
		Email:    dbUser.Email,
		Password: dbUser.PasswordHash,
	}
	return user, nil
}

func (r *AuthStore) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	dbUser, err := r.Queries.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User{
		ID:    dbUser.ID,
		Name:  dbUser.Name,
		Email: dbUser.Email,
	}
	return user, nil
}

func (r *AuthStore) UpdateRefreshToken(ctx context.Context, id uuid.UUID, refreshToken *string) error {
	return r.Queries.UpdateRefreshToken(ctx, generated.UpdateRefreshTokenParams{
		ID:                  id,
		CryptedRefreshToken: refreshToken,
	})
}
