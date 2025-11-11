package store

import (
	"context"
	"messenger-app/internal/auth/domain"
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

func (r *AuthRepository) CreateUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
	ref := generated.CreateUserParams(params)

	dbUser, err := r.Queries.CreateUser(ctx, ref)
	if err != nil {
		return domain.User{}, err
	}
	user := domain.User(dbUser)

	return user, nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	dbUser, err := r.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User(dbUser)

	return user, nil
}

func (r *AuthRepository) LogIn(ctx context.Context, params domain.LogInParams) (domain.User, error) {
	ref := generated.LogInParams(params)

	dbUser, err := r.Queries.LogIn(ctx, ref)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User(dbUser)

	return user, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, uuid uuid.UUID) (domain.User, error) {
	dbUser, err := r.Queries.GetUser(ctx, uuid)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User(dbUser)
	return user, nil
}

func (r *AuthRepository) UpdateRefreshToken(ctx context.Context, uuid uuid.UUID, refreshToken *string) error {
	if err := r.Queries.UpdateRefreshToken(ctx, generated.UpdateRefreshTokenParams{
		ID:                  uuid,
		CryptedRefreshToken: refreshToken,
	}); err != nil {
		return err
	}
	return nil
}
