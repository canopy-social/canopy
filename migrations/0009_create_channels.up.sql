
CREATE TABLE channels (
    id              TEXT PRIMARY KEY,
    account_id      TEXT UNIQUE NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    description     TEXT,
    description_html TEXT,
    rules           TEXT,
    owner_id        TEXT NOT NULL REFERENCES accounts(id),
    is_open         BOOLEAN NOT NULL DEFAULT TRUE,
    requires_approval BOOLEAN NOT NULL DEFAULT FALSE,
    topic           TEXT,
    members_count   INTEGER NOT NULL DEFAULT 0,
    posts_count     INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE channel_members (
    id              TEXT PRIMARY KEY,
    channel_id      TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    account_id      TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    role            TEXT NOT NULL DEFAULT 'member',
    status          TEXT NOT NULL DEFAULT 'active',
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (channel_id, account_id)
);

CREATE TABLE channel_invites (
    id              TEXT PRIMARY KEY,
    channel_id      TEXT NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    created_by      TEXT NOT NULL REFERENCES accounts(id),
    token           TEXT UNIQUE NOT NULL,
    invite_role     TEXT NOT NULL DEFAULT 'member',
    max_uses        INTEGER,
    uses_count      INTEGER NOT NULL DEFAULT 0,
    expires_at      TIMESTAMPTZ,
    is_revoked      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
