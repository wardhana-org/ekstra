DROP INDEX IF EXISTS user_auth_providers_unique_external_account;
DROP INDEX IF EXISTS user_auth_providers_one_password_per_user;

ALTER TABLE user_auth_providers
    ADD CONSTRAINT user_auth_providers_provider_provider_user_id_key
    UNIQUE (provider, provider_user_id);
