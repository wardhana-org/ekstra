package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wardhana-org/ekstra/backend/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	const query = `
		SELECT id, email, username, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindPasswordAuthProviderByUserID(ctx context.Context, userID int64) (*models.UserAuthProvider, error) {
	const query = `
		SELECT id, user_id, provider, provider_user_id, password_hash, created_at, updated_at
		FROM user_auth_providers
		WHERE user_id = $1 AND provider = 'password'
	`

	var provider models.UserAuthProvider

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&provider.ID,
		&provider.UserID,
		&provider.Provider,
		&provider.ProviderUserID,
		&provider.PasswordHash,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &provider, nil
}

type CredentialAvailability struct {
	EmailAvailable    bool
	UsernameAvailable bool
}

type CreatePasswordUserInput struct {
	Email        string
	Username     string
	PasswordHash string
}

func (r *UserRepository) CheckCredentialAvailability(ctx context.Context, email string, username string) (*CredentialAvailability, error) {
	var emailExists bool
	var usernameExists bool

	const availabilityCheckQuery = `
		SELECT
			EXISTS(SELECT 1 FROM users WHERE email = $1),
			EXISTS(SELECT 1 FROM users WHERE username = $2)
	`
	err := r.db.QueryRow(ctx, availabilityCheckQuery, email, username).Scan(&emailExists, &usernameExists)
	if err != nil {
		return nil, err
	}

	return &CredentialAvailability{
		EmailAvailable:    !emailExists,
		UsernameAvailable: !usernameExists,
	}, nil
}

func (r *UserRepository) CreatePasswordUser(ctx context.Context, input CreatePasswordUserInput) (*models.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const userQuery = `
		INSERT INTO users (email, username)
		VALUES ($1, $2)
		RETURNING id, email, username, status, created_at, updated_at
	`

	var user models.User

	err = tx.QueryRow(ctx, userQuery, input.Email, input.Username).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, mapUserConstraintError(err)
	}

	const authProviderQuery = `
		INSERT INTO user_auth_providers (user_id, provider, password_hash)
		VALUES ($1, 'password', $2)
	`

	if _, err = tx.Exec(ctx, authProviderQuery, user.ID, input.PasswordHash); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func mapUserConstraintError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return err
	}

	switch pgErr.ConstraintName {
	case "users_email_key":
		return ErrDuplicateEmail
	case "users_username_key":
		return ErrDuplicateUsername
	default:
		return err
	}
}
