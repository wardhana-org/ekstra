package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

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

func minTime(a time.Time, b time.Time) time.Time {
	if a.After(b) {
		return b
	}

	return a
}
