
CREATE TABLE page_themes (
    id              TEXT PRIMARY KEY,
    account_id      TEXT UNIQUE NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    colors          JSONB NOT NULL DEFAULT '{"background":"#ffffff","surface":"#f8f8f8","text_primary":"#111111","text_secondary":"#555555","accent":"#0066ff","accent_text":"#ffffff","border":"#e0e0e0","link":"#0066ff"}',
    fonts           JSONB NOT NULL DEFAULT '{"body":{"family":"system-ui","size":16,"weight":400,"line_height":1.6,"letter_spacing":0},"heading":{"family":"system-ui","size":24,"weight":700,"line_height":1.2,"letter_spacing":-0.02},"mono":{"family":"monospace","size":14,"weight":400,"line_height":1.5,"letter_spacing":0},"display":{"family":"system-ui","size":48,"weight":900,"line_height":1.0,"letter_spacing":-0.03}}',
    layout          JSONB NOT NULL DEFAULT '[]',
    stickers        JSONB NOT NULL DEFAULT '[]',
    widgets         JSONB NOT NULL DEFAULT '[]',
    bg_type         TEXT NOT NULL DEFAULT 'color',
    bg_gradient     JSONB,
    bg_image_id     TEXT REFERENCES media_attachments(id),
    bg_image_size   TEXT NOT NULL DEFAULT 'cover',
    bg_blur         SMALLINT NOT NULL DEFAULT 0,
    bg_opacity      SMALLINT NOT NULL DEFAULT 100,
    page_max_width  SMALLINT NOT NULL DEFAULT 800,
    page_padding    SMALLINT NOT NULL DEFAULT 24,
    show_follower_count BOOLEAN NOT NULL DEFAULT TRUE,
    show_following_count BOOLEAN NOT NULL DEFAULT TRUE,
    garden_mode     BOOLEAN NOT NULL DEFAULT FALSE,
    inherits_server_theme BOOLEAN NOT NULL DEFAULT FALSE,
    parent_theme_id TEXT REFERENCES page_themes(id),
    generated_css   TEXT,
    css_generated_at TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE theme_versions (
    id              TEXT PRIMARY KEY,
    account_id      TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    theme_snapshot  JSONB NOT NULL,
    label           TEXT,
    auto_saved      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX theme_versions_account ON theme_versions (account_id, created_at DESC);

CREATE TABLE essay_themes (
    id              TEXT PRIMARY KEY,
    essay_id        TEXT UNIQUE NOT NULL REFERENCES essays(id) ON DELETE CASCADE,
    colors          JSONB NOT NULL DEFAULT '{}',
    fonts           JSONB NOT NULL DEFAULT '{}',
    layout          JSONB NOT NULL DEFAULT '[]',
    bg_type         TEXT NOT NULL DEFAULT 'inherit',
    bg_gradient     JSONB,
    bg_image_id     TEXT REFERENCES media_attachments(id),
    generated_css   TEXT,
    css_generated_at TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE server_theme (
    id              TEXT PRIMARY KEY DEFAULT 'singleton',
    colors          JSONB NOT NULL DEFAULT '{}',
    fonts           JSONB NOT NULL DEFAULT '{}',
    layout          JSONB NOT NULL DEFAULT '[]',
    bg_type         TEXT NOT NULL DEFAULT 'color',
    generated_css   TEXT,
    css_generated_at TIMESTAMPTZ,
    updated_by      TEXT REFERENCES accounts(id),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT server_theme_singleton CHECK (id = 'singleton')
);
