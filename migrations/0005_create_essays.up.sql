-- 0005_create_essays.up.sql

CREATE TABLE essays (
    id              TEXT PRIMARY KEY,
    uri             TEXT UNIQUE NOT NULL,
    url             TEXT,
    account_id      TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    slug            TEXT NOT NULL,
    subtitle        TEXT,
    content         TEXT NOT NULL,
    content_text    TEXT NOT NULL,
    content_raw     TEXT NOT NULL,
    cover_media_id  TEXT REFERENCES media_attachments(id),
    reading_time_minutes INTEGER,
    visibility      TEXT NOT NULL DEFAULT 'public',
    language        TEXT,
    is_local        BOOLEAN NOT NULL DEFAULT TRUE,
    word_count      INTEGER NOT NULL DEFAULT 0,
    likes_count     INTEGER NOT NULL DEFAULT 0,
    views_count     INTEGER NOT NULL DEFAULT 0,
    published_at    TIMESTAMPTZ,
    ap_published    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (account_id, slug)
);

CREATE INDEX essays_account_id ON essays (account_id, published_at DESC);
CREATE INDEX essays_uri ON essays (uri);

ALTER TABLE essays ADD COLUMN content_tsv TSVECTOR
    GENERATED ALWAYS AS (to_tsvector('english', content_text || ' ' || title)) STORED;
CREATE INDEX essays_content_tsv ON essays USING GIN (content_tsv);

-- Essay marginalia
CREATE TABLE essay_marginalia (
    id              TEXT PRIMARY KEY,
    essay_id        TEXT NOT NULL REFERENCES essays(id) ON DELETE CASCADE,
    account_id      TEXT NOT NULL REFERENCES accounts(id),
    anchor_start    INTEGER NOT NULL,
    anchor_end      INTEGER NOT NULL,
    anchor_text     TEXT NOT NULL,
    content         TEXT NOT NULL,
    is_author_note  BOOLEAN NOT NULL DEFAULT FALSE,
    is_public       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX marginalia_essay ON essay_marginalia (essay_id, anchor_start);
CREATE INDEX marginalia_author ON essay_marginalia (essay_id) WHERE is_author_note = TRUE;

-- Add FK for media_attachments
ALTER TABLE media_attachments
    ADD CONSTRAINT fk_media_essay FOREIGN KEY (essay_id) REFERENCES essays(id);
