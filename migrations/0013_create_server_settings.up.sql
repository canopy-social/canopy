-- 0013_create_server_settings.up.sql

CREATE TABLE server_settings (
    key     TEXT PRIMARY KEY,
    value   JSONB NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE server_invites (
    id          TEXT PRIMARY KEY,
    created_by  TEXT NOT NULL REFERENCES accounts(id),
    token       TEXT UNIQUE NOT NULL,
    max_uses    INTEGER,
    uses_count  INTEGER NOT NULL DEFAULT 0,
    expires_at  TIMESTAMPTZ,
    is_revoked  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Garden mode positions
CREATE TABLE garden_positions (
    id              TEXT PRIMARY KEY,
    account_id      TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    item_type       TEXT NOT NULL,
    item_id         TEXT NOT NULL,
    x               FLOAT NOT NULL,
    y               FLOAT NOT NULL,
    width           FLOAT NOT NULL DEFAULT 20,
    z_index         SMALLINT NOT NULL DEFAULT 1,
    rotation        FLOAT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (account_id, item_type, item_id)
);
