-- 0002_create_media_attachments.up.sql
-- Media must be created before posts/essays since they reference it.

CREATE TABLE media_attachments (
    id              TEXT PRIMARY KEY,
    account_id      TEXT NOT NULL REFERENCES accounts(id),
    post_id         TEXT,
    essay_id        TEXT,
    remote_url      TEXT,
    url             TEXT,
    thumbnail_url   TEXT,
    type            TEXT NOT NULL,
    mime_type       TEXT,
    file_size       BIGINT,
    width           INTEGER,
    height          INTEGER,
    duration_seconds FLOAT,
    blurhash        TEXT,
    alt_text        TEXT,
    is_processed    BOOLEAN NOT NULL DEFAULT FALSE,
    processing_error TEXT,
    storage_key     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX media_attachments_account ON media_attachments (account_id);
CREATE INDEX media_attachments_post ON media_attachments (post_id) WHERE post_id IS NOT NULL;
CREATE INDEX media_attachments_essay ON media_attachments (essay_id) WHERE essay_id IS NOT NULL;
