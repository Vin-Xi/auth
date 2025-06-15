BEGIN;
CREATE TABLE IF NOT EXISTS auth.refresh_tokens (
    id              UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    token_hash      VARCHAR(255) UNIQUE NOT NULL,
    family_id       UUID NOT NULL,
    expires_at      TIMESTAMPTZ NOT NULL,
    revoked_at      TIMESTAMPTZ,
    created_ip      VARCHAR(45),
    user_agent      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON auth.refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_family_id ON auth.refresh_tokens(family_id);
COMMIT;