
CREATE TABLE notifications (
    id              TEXT PRIMARY KEY,
    account_id      TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    type            TEXT NOT NULL,
    from_account_id TEXT REFERENCES accounts(id),
    post_id         TEXT REFERENCES posts(id),
    essay_id        TEXT REFERENCES essays(id),
    channel_id      TEXT REFERENCES channels(id),
    read_at         TIMESTAMPTZ,
    dismissed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX notifications_account ON notifications (account_id, created_at DESC);
CREATE INDEX notifications_unread ON notifications (account_id) WHERE read_at IS NULL;
CREATE INDEX notifications_account_unread ON notifications (account_id, created_at DESC) WHERE read_at IS NULL;
