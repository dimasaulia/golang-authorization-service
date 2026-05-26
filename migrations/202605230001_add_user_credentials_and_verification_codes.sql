-- +goose Up
ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMPTZ NULL;

CREATE TABLE user_credentials (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    password_hash TEXT NOT NULL,
    must_change_password BOOLEAN NOT NULL DEFAULT FALSE,
    password_changed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT uq_user_credentials_user_id UNIQUE (user_id),
    CONSTRAINT fk_user_credentials_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_user_credentials_user_id ON user_credentials(user_id);
CREATE INDEX idx_user_credentials_must_change_password ON user_credentials(must_change_password);

CREATE TABLE user_verification_codes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    purpose VARCHAR(40) NOT NULL CHECK (purpose IN ('email_verification', 'password_setup', 'password_reset')),
    code_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NULL DEFAULT NOW(),
    CONSTRAINT fk_user_verification_codes_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_user_verification_codes_user_purpose ON user_verification_codes(user_id, purpose);
CREATE INDEX idx_user_verification_codes_code_hash ON user_verification_codes(code_hash);
CREATE INDEX idx_user_verification_codes_expires_at ON user_verification_codes(expires_at);

-- +goose Down
DROP TABLE IF EXISTS user_verification_codes;
DROP TABLE IF EXISTS user_credentials;
ALTER TABLE users DROP COLUMN IF EXISTS email_verified_at;
