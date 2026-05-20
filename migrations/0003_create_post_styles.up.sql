-- 0003_create_post_styles.up.sql
-- Per-post custom styling (separate table to keep posts clean).

CREATE TABLE post_styles (
    id              TEXT PRIMARY KEY,
    account_id      TEXT NOT NULL REFERENCES accounts(id),
    background_color    TEXT,
    background_image_id TEXT REFERENCES media_attachments(id),
    text_color          TEXT,
    font_family         TEXT,
    font_size           SMALLINT,
    font_weight         SMALLINT,
    border_radius       SMALLINT,
    border_color        TEXT,
    border_width        SMALLINT,
    padding             SMALLINT,
    has_texture         BOOLEAN NOT NULL DEFAULT FALSE,
    texture_type        TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
