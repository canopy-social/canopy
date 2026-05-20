-- 0004_create_posts.up.sql

CREATE TABLE posts (
    id              TEXT PRIMARY KEY,
    uri             TEXT UNIQUE NOT NULL,
    url             TEXT,
    account_id      TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    content         TEXT NOT NULL,
    content_text    TEXT NOT NULL,
    content_warning TEXT,
    is_sensitive    BOOLEAN NOT NULL DEFAULT FALSE,
    visibility      TEXT NOT NULL DEFAULT 'public',
    language        TEXT,
    reply_to_id     TEXT REFERENCES posts(id),
    reply_to_uri    TEXT,
    thread_root_id  TEXT REFERENCES posts(id),
    boost_of_id     TEXT REFERENCES posts(id),
    boost_of_uri    TEXT,
    is_local        BOOLEAN NOT NULL DEFAULT TRUE,
    is_pinned       BOOLEAN NOT NULL DEFAULT FALSE,
    ap_published    TIMESTAMPTZ,
    edit_history    JSONB NOT NULL DEFAULT '[]',
    likes_count     INTEGER NOT NULL DEFAULT 0,
    boosts_count    INTEGER NOT NULL DEFAULT 0,
    replies_count   INTEGER NOT NULL DEFAULT 0,
    post_style_id   TEXT REFERENCES post_styles(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX posts_account_id ON posts (account_id, created_at DESC);
CREATE INDEX posts_uri ON posts (uri);
CREATE INDEX posts_thread_root ON posts (thread_root_id) WHERE thread_root_id IS NOT NULL;
CREATE INDEX posts_visibility ON posts (visibility, created_at DESC);
CREATE INDEX posts_visibility_public ON posts (visibility, created_at DESC) WHERE visibility = 'public';
CREATE INDEX posts_reply_to ON posts (reply_to_id) WHERE reply_to_id IS NOT NULL;

ALTER TABLE posts ADD COLUMN content_tsv TSVECTOR
    GENERATED ALWAYS AS (to_tsvector('english', content_text)) STORED;
CREATE INDEX posts_content_tsv ON posts USING GIN (content_tsv);

-- Post mentions join table
CREATE TABLE post_mentions (
    post_id     TEXT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    account_id  TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    uri         TEXT,
    PRIMARY KEY (post_id, account_id)
);

-- Post tags (hashtags) join table
CREATE TABLE post_tags (
    post_id     TEXT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag         TEXT NOT NULL,
    PRIMARY KEY (post_id, tag)
);

CREATE INDEX post_tags_tag ON post_tags (tag);

-- Likes
CREATE TABLE post_likes (
    id          TEXT PRIMARY KEY,
    post_id     TEXT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    account_id  TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    uri         TEXT UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (post_id, account_id)
);

-- Boosts (stored as separate records, not as posts with boost_of_id for count tracking)
CREATE TABLE post_boosts (
    id          TEXT PRIMARY KEY,
    post_id     TEXT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    account_id  TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    uri         TEXT UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (post_id, account_id)
);

-- Add FKs for media_attachments now that posts exist
ALTER TABLE media_attachments
    ADD CONSTRAINT fk_media_post FOREIGN KEY (post_id) REFERENCES posts(id);
