package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/wardhana-org/ekstra/backend/internal/models"
)

type CreateSessionInput struct {
	UserID            int64
	ClientType        string
	DeviceName        *string
	UserAgent         *string
	IPAddress         *string
	ExpiresAt         time.Time
	AbsoluteExpiresAt time.Time
}

type CreateTokenInput struct {
	TokenHash  string
	TokenType  string
	ExpiresAt  time.Time
	Generation *int
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

	absoluteExpiresAt := input.Session.AbsoluteExpiresAt
	if absoluteExpiresAt.IsZero() {
		absoluteExpiresAt = input.Session.ExpiresAt
	}

	const sessionQuery = `
		INSERT INTO auth_sessions (
			user_id,
			client_type,
			device_name,
			user_agent,
			ip_address,
			expires_at,
			absolute_expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, client_type, device_name, user_agent, ip_address, created_at, last_seen_at, expires_at, absolute_expires_at, revoked_at, revoked_reason, reuse_detected_at
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
		absoluteExpiresAt,
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
		&session.AbsoluteExpiresAt,
		&session.RevokedAt,
		&session.RevokedReason,
		&session.ReuseDetectedAt,
	)
	if err != nil {
		return nil, nil, err
	}

	const tokenQuery = `
		INSERT INTO auth_tokens (
			session_id,
			token_hash,
			token_type,
			expires_at,
			generation
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, session_id, token_hash, token_type, generation, created_at, expires_at, used_at, replaced_by_token_id, reuse_grace_expires_at, revoked_at, revoked_reason
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
			tokenInput.Generation,
		).Scan(
			&token.ID,
			&token.SessionID,
			&token.TokenHash,
			&token.TokenType,
			&token.Generation,
			&token.CreatedAt,
			&token.ExpiresAt,
			&token.UsedAt,
			&token.ReplacedByTokenID,
			&token.ReuseGraceExpiresAt,
			&token.RevokedAt,
			&token.RevokedReason,
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
			AND s.absolute_expires_at > $2
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

	if err = revokeSessionTokens(
		ctx,
		tx,
		sessionID,
		revokedAt,
		authSessionRevokedReasonLogout,
		authTokenRevokedReasonLogout,
		false,
	); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func revokeSessionTokens(
	ctx context.Context,
	tx pgx.Tx,
	sessionID int64,
	revokedAt time.Time,
	sessionRevokedReason string,
	tokenRevokedReason string,
	reuseDetected bool,
) error {
	const revokeSessionQuery = `
		UPDATE auth_sessions
		SET revoked_at = $2,
			revoked_reason = $3,
			reuse_detected_at = CASE
				WHEN $4 THEN $2
				ELSE reuse_detected_at
			END
		WHERE id = $1
			AND revoked_at IS NULL
	`

	if _, err := tx.Exec(ctx, revokeSessionQuery, sessionID, revokedAt, sessionRevokedReason, reuseDetected); err != nil {
		return err
	}

	const revokeTokensQuery = `
		UPDATE auth_tokens
		SET revoked_at = $2,
			revoked_reason = $3
		WHERE session_id = $1
			AND revoked_at IS NULL
	`

	if _, err := tx.Exec(ctx, revokeTokensQuery, sessionID, revokedAt, tokenRevokedReason); err != nil {
		return err
	}

	return nil
}
