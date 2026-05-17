package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
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

type RotateRefreshTokenInput struct {
	RefreshTokenHash      string
	NewAccessTokenHash    string
	NewRefreshTokenHash   string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	RotatedAt             time.Time
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

func (r *AuthRepository) FindUserByAccessTokenHash(ctx context.Context, tokenHash string, now time.Time) (*models.User, error) {
	const query = `
		SELECT u.id, u.email, u.username, u.status, u.created_at, u.updated_at
		FROM auth_tokens AS t
		JOIN auth_sessions AS s ON s.id = t.session_id
		JOIN users AS u ON u.id = s.user_id
		WHERE t.token_hash = $1
			AND t.token_type = 'access'
			AND t.expires_at > $2
			AND t.revoked_at IS NULL
			AND s.expires_at > $2
			AND s.revoked_at IS NULL
	`

	var user models.User

	err := r.db.QueryRow(ctx, query, tokenHash, now).Scan(
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

func (r *AuthRepository) RevokeSessionByRefreshTokenHash(ctx context.Context, tokenHash string, revokedAt time.Time) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const sessionQuery = `
		SELECT session_id
		FROM auth_tokens
		WHERE token_hash = $1
			AND token_type = 'refresh'
	`

	var sessionID int64

	err = tx.QueryRow(ctx, sessionQuery, tokenHash).Scan(&sessionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}

		return err
	}

	const revokeSessionQuery = `
		UPDATE auth_sessions
		SET revoked_at = $2
		WHERE id = $1
			AND revoked_at IS NULL
	`

	if _, err = tx.Exec(ctx, revokeSessionQuery, sessionID, revokedAt); err != nil {
		return err
	}

	const revokeTokensQuery = `
		UPDATE auth_tokens
		SET revoked_at = $2
		WHERE session_id = $1
			AND revoked_at IS NULL
	`

	if _, err = tx.Exec(ctx, revokeTokensQuery, sessionID, revokedAt); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) RotateRefreshToken(ctx context.Context, input RotateRefreshTokenInput) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const sessionQuery = `
		SELECT t.id, s.id
		FROM auth_tokens AS t
		JOIN auth_sessions AS s ON s.id = t.session_id
		WHERE t.token_hash = $1
			AND t.token_type = 'refresh'
			AND t.expires_at > $2
			AND t.revoked_at IS NULL
			AND s.expires_at > $2
			AND s.revoked_at IS NULL
		FOR UPDATE OF t, s
	`

	var refreshTokenID int64
	var sessionID int64

	err = tx.QueryRow(ctx, sessionQuery, input.RefreshTokenHash, input.RotatedAt).Scan(
		&refreshTokenID,
		&sessionID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}

		return err
	}

	const revokeRefreshTokenQuery = `
		UPDATE auth_tokens
		SET revoked_at = $2
		WHERE id = $1
	`

	if _, err = tx.Exec(ctx, revokeRefreshTokenQuery, refreshTokenID, input.RotatedAt); err != nil {
		return err
	}

	const updateSessionQuery = `
		UPDATE auth_sessions
		SET last_seen_at = $2,
			expires_at = $3
		WHERE id = $1
	`

	if _, err = tx.Exec(ctx, updateSessionQuery, sessionID, input.RotatedAt, input.RefreshTokenExpiresAt); err != nil {
		return err
	}

	const tokenQuery = `
		INSERT INTO auth_tokens (
			session_id,
			token_hash,
			token_type,
			expires_at
		)
		VALUES ($1, $2, $3, $4)
	`

	if _, err = tx.Exec(ctx, tokenQuery, sessionID, input.NewAccessTokenHash, "access", input.AccessTokenExpiresAt); err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, tokenQuery, sessionID, input.NewRefreshTokenHash, "refresh", input.RefreshTokenExpiresAt); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
