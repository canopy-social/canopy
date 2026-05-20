-- 0012_create_moderation.up.sql

CREATE TABLE reports (
    id              TEXT PRIMARY KEY,
    reporter_id     TEXT NOT NULL REFERENCES accounts(id),
    target_account_id TEXT REFERENCES accounts(id),
    target_post_id  TEXT REFERENCES posts(id),
    target_essay_id TEXT REFERENCES essays(id),
    category        TEXT NOT NULL,
    comment         TEXT,
    status          TEXT NOT NULL DEFAULT 'open',
    resolved_by     TEXT REFERENCES accounts(id),
    resolved_at     TIMESTAMPTZ,
    action_taken    TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX reports_status ON reports (status, created_at DESC);

CREATE TABLE instance_blocks (
    id              TEXT PRIMARY KEY,
    domain          TEXT UNIQUE NOT NULL,
    severity        TEXT NOT NULL DEFAULT 'silence',
    reason          TEXT,
    reject_media    BOOLEAN NOT NULL DEFAULT FALSE,
    reject_reports  BOOLEAN NOT NULL DEFAULT FALSE,
    created_by      TEXT REFERENCES accounts(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
