package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wardhana-org/ekstra/backend/internal/models"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

type CreateSessionInput struct {
	UserID     int64
	ClientType string
	DeviceName *string
	UserAgent  *string
	IPAddress  *string
	ExpiresAt  time.Time
}

type CreateTokenInput struct {
	TokenHash string
	TokenType string
	ExpiresAt time.Time
}

type CreateSessionWithTokensInput struct {
	Session CreateSessionInput
	Tokens  []CreateTokenInput
}

func (r *AuthRepository) CreateSessionWithTokens(ctx context.Context, input CreateSessionWithTokensInput) (*models.AuthSession, []models.AuthToken, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	const sessionQuery = `
		INSERT INTO auth_sessions (
			user_id,
			client_type,
			device_name,
			user_agent,
			ip_address,
			expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, client_type, device_name, user_agent, ip_address, created_at, last_seen_at, expires_at, revoked_at
	`

	var session models.AuthSession

	err = tx.QueryRow(
		ctx,
		sessionQuery,
		input.Session.UserID,
		input.Session.ClientType,
		input.Session.DeviceName,
		input.Session.UserAgent,
		input.Session.IPAddress,
		input.Session.ExpiresAt,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.ClientType,
		&session.DeviceName,
		&session.UserAgent,
		&session.IPAddress,
		&session.CreatedAt,
		&session.LastSeenAt,
		&session.ExpiresAt,
		&session.RevokedAt,
	)
	if err != nil {
		return nil, nil, err
	}

	const tokenQuery = `
		INSERT INTO auth_tokens (
			session_id,
			token_hash,
			token_type,
			expires_at
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, session_id, token_hash, token_type, created_at, expires_at, revoked_at
	`

	tokens := make([]models.AuthToken, 0, len(input.Tokens))
	for _, tokenInput := range input.Tokens {
		var token models.AuthToken

		err = tx.QueryRow(
			ctx,
			tokenQuery,
			session.ID,
			tokenInput.TokenHash,
			tokenInput.TokenType,
			tokenInput.ExpiresAt,
		).Scan(
			&token.ID,
			&token.SessionID,
			&token.TokenHash,
			&token.TokenType,
			&token.CreatedAt,
			&token.ExpiresAt,
			&token.RevokedAt,
		)
		if err != nil {
			return nil, nil, err
		}

		tokens = append(tokens, token)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return &session, tokens, nil
}
