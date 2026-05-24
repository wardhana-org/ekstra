ALTER TABLE auth_sessions
    ADD COLUMN IF NOT EXISTS absolute_expires_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS revoked_reason VARCHAR(100),
    ADD COLUMN IF NOT EXISTS reuse_detected_at TIMESTAMPTZ;

UPDATE auth_sessions
SET absolute_expires_at = created_at + INTERVAL '90 days'
WHERE absolute_expires_at IS NULL;

ALTER TABLE auth_sessions
    ALTER COLUMN absolute_expires_at SET NOT NULL;

ALTER TABLE auth_tokens
    ADD COLUMN IF NOT EXISTS generation INTEGER,
    ADD COLUMN IF NOT EXISTS used_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS replaced_by_token_id BIGINT REFERENCES auth_tokens(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS reuse_grace_expires_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS revoked_reason VARCHAR(100);

UPDATE auth_tokens
SET generation = 1
WHERE token_type = 'refresh'
    AND generation IS NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'chk_auth_tokens_refresh_generation'
    ) THEN
        ALTER TABLE auth_tokens
            ADD CONSTRAINT chk_auth_tokens_refresh_generation
            CHECK (token_type <> 'refresh' OR generation IS NOT NULL);
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'chk_auth_tokens_generation_positive'
    ) THEN
        ALTER TABLE auth_tokens
            ADD CONSTRAINT chk_auth_tokens_generation_positive
            CHECK (generation IS NULL OR generation > 0);
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS auth_security_events (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    session_id BIGINT REFERENCES auth_sessions(id) ON DELETE SET NULL,
    token_id BIGINT REFERENCES auth_tokens(id) ON DELETE SET NULL,
    event_type VARCHAR(100) NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_auth_tokens_refresh_chain ON auth_tokens(session_id, token_type, generation)
    WHERE token_type = 'refresh';
CREATE INDEX IF NOT EXISTS idx_auth_tokens_replaced_by ON auth_tokens(replaced_by_token_id)
    WHERE replaced_by_token_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_auth_security_events_user_id ON auth_security_events(user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_auth_security_events_session_id ON auth_security_events(session_id, created_at);
CREATE INDEX IF NOT EXISTS idx_auth_security_events_event_type ON auth_security_events(event_type, created_at);
