ALTER TABLE user_auth_providers
    DROP CONSTRAINT IF EXISTS user_auth_providers_provider_provider_user_id_key;

CREATE UNIQUE INDEX IF NOT EXISTS user_auth_providers_one_password_per_user
    ON user_auth_providers (user_id)
    WHERE provider = 'password';

CREATE UNIQUE INDEX IF NOT EXISTS user_auth_providers_unique_external_account
    ON user_auth_providers (provider, provider_user_id)
    WHERE provider_user_id IS NOT NULL;
