DROP INDEX IF EXISTS idx_auth_security_events_event_type;
DROP INDEX IF EXISTS idx_auth_security_events_session_id;
DROP INDEX IF EXISTS idx_auth_security_events_user_id;
DROP INDEX IF EXISTS idx_auth_tokens_replaced_by;
DROP INDEX IF EXISTS idx_auth_tokens_refresh_chain;

DROP TABLE IF EXISTS auth_security_events;

ALTER TABLE auth_tokens
    DROP CONSTRAINT IF EXISTS chk_auth_tokens_generation_positive,
    DROP CONSTRAINT IF EXISTS chk_auth_tokens_refresh_generation;

ALTER TABLE auth_tokens
    DROP COLUMN IF EXISTS revoked_reason,
    DROP COLUMN IF EXISTS reuse_grace_expires_at,
    DROP COLUMN IF EXISTS replaced_by_token_id,
    DROP COLUMN IF EXISTS used_at,
    DROP COLUMN IF EXISTS generation;

ALTER TABLE auth_sessions
    DROP COLUMN IF EXISTS reuse_detected_at,
    DROP COLUMN IF EXISTS revoked_reason,
    DROP COLUMN IF EXISTS absolute_expires_at;
