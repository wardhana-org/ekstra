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

const (
	authTokenTypeAccess  = "access"
	authTokenTypeRefresh = "refresh"

	authSessionRevokedReasonLogout            = "logout"
	authSessionRevokedReasonRefreshTokenReuse = "refresh_token_reuse"

	authTokenRevokedReasonLogout            = "logout"
	authTokenRevokedReasonRotated           = "rotated"
	authTokenRevokedReasonRefreshTokenReuse = "refresh_token_reuse"

	authSecurityEventRefreshTokenRace   = "refresh_token_race"
	authSecurityEventRefreshTokenReused = "refresh_token_reused"
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

type RotateRefreshTokenInput struct {
	RefreshTokenHash      string
	NewAccessTokenHash    string
	NewRefreshTokenHash   string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	RotatedAt             time.Time
	ReuseGracePeriod      time.Duration
	UserAgent             *string
	IPAddress             *string
}

type RotateRefreshTokenResult struct {
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
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

func (r *AuthRepository) RotateRefreshToken(ctx context.Context, input RotateRefreshTokenInput) (*RotateRefreshTokenResult, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const sessionQuery = `
		SELECT t.id, s.id, COALESCE(t.generation, 1), s.absolute_expires_at
		FROM auth_tokens AS t
		JOIN auth_sessions AS s ON s.id = t.session_id
		WHERE t.token_hash = $1
			AND t.token_type = 'refresh'
			AND t.expires_at > $2
			AND t.revoked_at IS NULL
			AND s.expires_at > $2
			AND s.absolute_expires_at > $2
			AND s.revoked_at IS NULL
		FOR UPDATE OF t, s
	`

	var refreshTokenID int64
	var sessionID int64
	var refreshTokenGeneration int
	var absoluteExpiresAt time.Time

	err = tx.QueryRow(ctx, sessionQuery, input.RefreshTokenHash, input.RotatedAt).Scan(
		&refreshTokenID,
		&sessionID,
		&refreshTokenGeneration,
		&absoluteExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			action, err := handleInvalidRefreshToken(ctx, tx, invalidRefreshTokenInput{
				TokenHash: input.RefreshTokenHash,
				Now:       input.RotatedAt,
				UserAgent: input.UserAgent,
				IPAddress: input.IPAddress,
			})
			if err != nil {
				return nil, err
			}
			if action == invalidRefreshTokenRace {
				if err := tx.Commit(ctx); err != nil {
					return nil, err
				}

				return nil, ErrRefreshTokenRace
			}
			if action == invalidRefreshTokenReuse {
				if err := tx.Commit(ctx); err != nil {
					return nil, err
				}

				return nil, ErrRefreshTokenReused
			}

			return nil, ErrNotFound
		}

		return nil, err
	}

	accessTokenExpiresAt := minTime(input.AccessTokenExpiresAt, absoluteExpiresAt)
	refreshTokenExpiresAt := minTime(input.RefreshTokenExpiresAt, absoluteExpiresAt)
	reuseGraceExpiresAt := input.RotatedAt.Add(input.ReuseGracePeriod)
	newRefreshGeneration := refreshTokenGeneration + 1

	const updateSessionQuery = `
		UPDATE auth_sessions
		SET last_seen_at = $2,
			expires_at = $3
		WHERE id = $1
	`

	if _, err = tx.Exec(ctx, updateSessionQuery, sessionID, input.RotatedAt, refreshTokenExpiresAt); err != nil {
		return nil, err
	}

	const accessTokenQuery = `
		INSERT INTO auth_tokens (
			session_id,
			token_hash,
			token_type,
			expires_at
		)
		VALUES ($1, $2, $3, $4)
	`

	if _, err = tx.Exec(ctx, accessTokenQuery, sessionID, input.NewAccessTokenHash, authTokenTypeAccess, accessTokenExpiresAt); err != nil {
		return nil, err
	}

	const refreshTokenQuery = `
		INSERT INTO auth_tokens (
			session_id,
			token_hash,
			token_type,
			expires_at,
			generation
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var newRefreshTokenID int64
	err = tx.QueryRow(
		ctx,
		refreshTokenQuery,
		sessionID,
		input.NewRefreshTokenHash,
		authTokenTypeRefresh,
		refreshTokenExpiresAt,
		newRefreshGeneration,
	).Scan(&newRefreshTokenID)
	if err != nil {
		return nil, err
	}

	const rotateRefreshTokenQuery = `
		UPDATE auth_tokens
		SET revoked_at = $2,
			used_at = $2,
			replaced_by_token_id = $3,
			reuse_grace_expires_at = $4,
			revoked_reason = $5
		WHERE id = $1
	`

	if _, err = tx.Exec(
		ctx,
		rotateRefreshTokenQuery,
		refreshTokenID,
		input.RotatedAt,
		newRefreshTokenID,
		reuseGraceExpiresAt,
		authTokenRevokedReasonRotated,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &RotateRefreshTokenResult{
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}

type invalidRefreshTokenAction int

const (
	invalidRefreshTokenUnknown invalidRefreshTokenAction = iota
	invalidRefreshTokenRace
	invalidRefreshTokenReuse
)

type invalidRefreshTokenInput struct {
	TokenHash string
	Now       time.Time
	UserAgent *string
	IPAddress *string
}

func handleInvalidRefreshToken(ctx context.Context, tx pgx.Tx, input invalidRefreshTokenInput) (invalidRefreshTokenAction, error) {
	const query = `
		SELECT
			t.id,
			t.session_id,
			t.revoked_at,
			t.revoked_reason,
			t.replaced_by_token_id,
			t.reuse_grace_expires_at,
			s.user_id,
			s.expires_at,
			s.absolute_expires_at,
			s.revoked_at,
			replacement.expires_at,
			replacement.revoked_at
		FROM auth_tokens AS t
		JOIN auth_sessions AS s ON s.id = t.session_id
		LEFT JOIN auth_tokens AS replacement ON replacement.id = t.replaced_by_token_id
		WHERE t.token_hash = $1
			AND t.token_type = 'refresh'
		FOR UPDATE OF t, s
	`

	var tokenID int64
	var sessionID int64
	var tokenRevokedAt *time.Time
	var tokenRevokedReason *string
	var replacedByTokenID *int64
	var reuseGraceExpiresAt *time.Time
	var userID int64
	var sessionExpiresAt time.Time
	var sessionAbsoluteExpiresAt time.Time
	var sessionRevokedAt *time.Time
	var replacementExpiresAt *time.Time
	var replacementRevokedAt *time.Time

	err := tx.QueryRow(ctx, query, input.TokenHash).Scan(
		&tokenID,
		&sessionID,
		&tokenRevokedAt,
		&tokenRevokedReason,
		&replacedByTokenID,
		&reuseGraceExpiresAt,
		&userID,
		&sessionExpiresAt,
		&sessionAbsoluteExpiresAt,
		&sessionRevokedAt,
		&replacementExpiresAt,
		&replacementRevokedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return invalidRefreshTokenUnknown, nil
		}

		return invalidRefreshTokenUnknown, err
	}

	if tokenRevokedAt == nil ||
		sessionRevokedAt != nil ||
		!sessionExpiresAt.After(input.Now) ||
		!sessionAbsoluteExpiresAt.After(input.Now) {
		return invalidRefreshTokenUnknown, nil
	}

	if isRefreshTokenRace(
		input.Now,
		tokenRevokedReason,
		replacedByTokenID,
		reuseGraceExpiresAt,
		replacementExpiresAt,
		replacementRevokedAt,
	) {
		if err := insertAuthSecurityEvent(ctx, tx, authSecurityEventInput{
			UserID:    &userID,
			SessionID: &sessionID,
			TokenID:   &tokenID,
			EventType: authSecurityEventRefreshTokenRace,
			IPAddress: input.IPAddress,
			UserAgent: input.UserAgent,
		}); err != nil {
			return invalidRefreshTokenUnknown, err
		}

		return invalidRefreshTokenRace, nil
	}

	if err := revokeSessionTokens(
		ctx,
		tx,
		sessionID,
		input.Now,
		authSessionRevokedReasonRefreshTokenReuse,
		authTokenRevokedReasonRefreshTokenReuse,
		true,
	); err != nil {
		return invalidRefreshTokenUnknown, err
	}

	if err := insertAuthSecurityEvent(ctx, tx, authSecurityEventInput{
		UserID:    &userID,
		SessionID: &sessionID,
		TokenID:   &tokenID,
		EventType: authSecurityEventRefreshTokenReused,
		IPAddress: input.IPAddress,
		UserAgent: input.UserAgent,
	}); err != nil {
		return invalidRefreshTokenUnknown, err
	}

	return invalidRefreshTokenReuse, nil
}

func isRefreshTokenRace(
	now time.Time,
	tokenRevokedReason *string,
	replacedByTokenID *int64,
	reuseGraceExpiresAt *time.Time,
	replacementExpiresAt *time.Time,
	replacementRevokedAt *time.Time,
) bool {
	if tokenRevokedReason == nil || *tokenRevokedReason != authTokenRevokedReasonRotated {
		return false
	}

	if replacedByTokenID == nil || reuseGraceExpiresAt == nil || !reuseGraceExpiresAt.After(now) {
		return false
	}

	if replacementExpiresAt == nil || !replacementExpiresAt.After(now) {
		return false
	}

	return replacementRevokedAt == nil
}

type authSecurityEventInput struct {
	UserID    *int64
	SessionID *int64
	TokenID   *int64
	EventType string
	IPAddress *string
	UserAgent *string
}

func insertAuthSecurityEvent(ctx context.Context, tx pgx.Tx, input authSecurityEventInput) error {
	const query = `
		INSERT INTO auth_security_events (
			user_id,
			session_id,
			token_id,
			event_type,
			ip_address,
			user_agent
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := tx.Exec(
		ctx,
		query,
		input.UserID,
		input.SessionID,
		input.TokenID,
		input.EventType,
		input.IPAddress,
		input.UserAgent,
	)

	return err
}

func minTime(a time.Time, b time.Time) time.Time {
	if a.After(b) {
		return b
	}

	return a
}
