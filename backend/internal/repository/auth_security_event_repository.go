package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

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
