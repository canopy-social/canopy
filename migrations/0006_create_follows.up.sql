
CREATE TABLE follows (
    id              TEXT PRIMARY KEY,
    follower_id     TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    following_id    TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    status          TEXT NOT NULL DEFAULT 'accepted',
    uri             TEXT UNIQUE,
    notify          BOOLEAN NOT NULL DEFAULT FALSE,
    show_boosts     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (follower_id, following_id)
);

CREATE INDEX follows_follower ON follows (follower_id, status);
CREATE INDEX follows_following ON follows (following_id, status);
CREATE INDEX follows_following_accepted ON follows (following_id) WHERE status = 'accepted';
CREATE INDEX follows_follower_accepted ON follows (follower_id) WHERE status = 'accepted';

CREATE TABLE blocks (
    id              TEXT PRIMARY KEY,
    account_id      TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    target_id       TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    uri             TEXT UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (account_id, target_id)
);

CREATE TABLE mutes (
    id              TEXT PRIMARY KEY,
    account_id      TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    target_id       TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    hide_notifications BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (account_id, target_id)
);
