package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type IUserRepository interface {
	GetById(ctx context.Context, userId uuid.UUID) (*User, error)
	GetByProviderId(ctx context.Context, provider string, providerId string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Add(ctx context.Context, params CreateUserParams) (*User, error)
}

type UserRepository struct {
	db *sql.DB
}

func CreateUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repository *UserRepository) GetById(ctx context.Context, userId uuid.UUID) (*User, error) {
	db := New(repository.db)
	user, err := db.GetUserById(ctx, userId)
	return &user, err
}

func (repository *UserRepository) GetByProviderId(ctx context.Context, provider string, providerId string) (*User, error) {
	db := New(repository.db)
	user, err := db.GetUserByProviderId(ctx, GetUserByProviderIdParams{
		Provider:   provider,
		ProviderID: providerId,
	})
	return &user, err
}

func (repository *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	db := New(repository.db)
	user, err := db.GetUserByEmail(ctx, email)
	return &user, err
}

func (repository *UserRepository) Add(ctx context.Context, params CreateUserParams) (*User, error) {
	db := New(repository.db)
	user, err := db.CreateUser(ctx, params)
	return &user, err
}
