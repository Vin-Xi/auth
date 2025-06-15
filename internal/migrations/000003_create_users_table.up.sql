BEGIN;
CREATE TABLE IF NOT EXISTS auth.users (
    id              UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    username        VARCHAR(50) UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    first_name      VARCHAR(50),
    last_name       VARCHAR(50),
    is_active       BOOLEAN NOT NULL DEFAULT true,
    -- is_verified     BOOLEAN NOT NULL DEFAULT false,
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);
CREATE INDEX IF NOT EXISTS idx_users_email ON auth.users(email);
COMMIT;