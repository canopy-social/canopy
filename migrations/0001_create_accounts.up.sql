-- 0001_create_accounts.up.sql
-- Core accounts table for both local and remote actors.

CREATE TABLE accounts (
    id                  TEXT PRIMARY KEY,
    username            TEXT NOT NULL,
    domain              TEXT,
    uri                 TEXT UNIQUE NOT NULL,
    display_name        TEXT,
    bio                 TEXT,
    bio_text            TEXT,
    avatar_media_id     TEXT,
    header_media_id     TEXT,
    public_key_pem      TEXT NOT NULL,
    private_key_pem     TEXT,
    key_id              TEXT NOT NULL,
    role                TEXT NOT NULL DEFAULT 'user',
    is_local            BOOLEAN NOT NULL DEFAULT FALSE,
    is_locked           BOOLEAN NOT NULL DEFAULT FALSE,
    is_bot              BOOLEAN NOT NULL DEFAULT FALSE,
    is_suspended        BOOLEAN NOT NULL DEFAULT FALSE,
    is_silenced         BOOLEAN NOT NULL DEFAULT FALSE,
    actor_type          TEXT NOT NULL DEFAULT 'Person',
    inbox_url           TEXT,
    outbox_url          TEXT,
    shared_inbox_url    TEXT,
    followers_url       TEXT,
    following_url       TEXT,
    featured_url        TEXT,
    followers_count     INTEGER NOT NULL DEFAULT 0,
    following_count     INTEGER NOT NULL DEFAULT 0,
    posts_count         INTEGER NOT NULL DEFAULT 0,
    last_fetched_at     TIMESTAMPTZ,
    custom_domain       TEXT UNIQUE,
    custom_domain_verified BOOLEAN NOT NULL DEFAULT FALSE,
    custom_domain_verification_token TEXT,
    password_hash       TEXT,
    email               TEXT UNIQUE,
    email_verified_at   TIMESTAMPTZ,
    email_verify_token  TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX accounts_username_domain ON accounts (username, domain);
CREATE INDEX accounts_uri ON accounts (uri);
CREATE INDEX accounts_custom_domain ON accounts (custom_domain) WHERE custom_domain IS NOT NULL;
CREATE INDEX accounts_is_local ON accounts (is_local);
CREATE INDEX accounts_email ON accounts (email) WHERE email IS NOT NULL;
