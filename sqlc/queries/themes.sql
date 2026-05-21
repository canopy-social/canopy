-- name: CreatePageTheme :one
INSERT INTO page_themes (id, account_id, colors, fonts, layout, stickers, widgets, bg_type, bg_gradient, bg_image_id, bg_image_size, bg_blur, bg_opacity, page_max_width, page_padding, show_follower_count, show_following_count, garden_mode, inherits_server_theme, parent_theme_id, generated_css, css_generated_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
RETURNING *;

-- name: GetPageThemeByID :one
SELECT * FROM page_themes WHERE id = $1;

-- name: GetPageThemeByAccountID :one
SELECT * FROM page_themes WHERE account_id = $1;

-- name: UpdatePageTheme :exec
UPDATE page_themes SET
  colors = $2, fonts = $3, layout = $4, stickers = $5, widgets = $6,
  bg_type = $7, bg_gradient = $8, bg_image_id = $9, bg_image_size = $10,
  bg_blur = $11, bg_opacity = $12, page_max_width = $13, page_padding = $14,
  show_follower_count = $15, show_following_count = $16, garden_mode = $17,
  inherits_server_theme = $18, parent_theme_id = $19, generated_css = $20,
  css_generated_at = $21, updated_at = $22
WHERE id = $1;

-- name: DeletePageTheme :exec
DELETE FROM page_themes WHERE id = $1;

-- name: CreateThemeVersion :one
INSERT INTO theme_versions (id, account_id, theme_snapshot, label, auto_saved, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListThemeVersions :many
SELECT * FROM theme_versions
WHERE account_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetThemeVersionByID :one
SELECT * FROM theme_versions WHERE id = $1;

-- name: DeleteThemeVersion :exec
DELETE FROM theme_versions WHERE id = $1;

-- name: CountAutoSavedVersions :one
SELECT COUNT(*)::int FROM theme_versions WHERE account_id = $1 AND auto_saved = TRUE;

-- name: CountNamedVersions :one
SELECT COUNT(*)::int FROM theme_versions WHERE account_id = $1 AND auto_saved = FALSE;

-- name: DeleteOldestAutoSave :exec
DELETE FROM theme_versions
WHERE id = (
  SELECT id FROM theme_versions
  WHERE account_id = $1 AND auto_saved = TRUE
  ORDER BY created_at ASC
  LIMIT 1
);

-- name: CreateEssayTheme :one
INSERT INTO essay_themes (id, essay_id, colors, fonts, layout, bg_type, bg_gradient, bg_image_id, generated_css, css_generated_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetEssayThemeByEssayID :one
SELECT * FROM essay_themes WHERE essay_id = $1;

-- name: UpdateEssayTheme :exec
UPDATE essay_themes SET
  colors = $2, fonts = $3, layout = $4, bg_type = $5, bg_gradient = $6,
  bg_image_id = $7, generated_css = $8, css_generated_at = $9, updated_at = $10
WHERE essay_id = $1;

-- name: DeleteEssayTheme :exec
DELETE FROM essay_themes WHERE essay_id = $1;

-- name: GetServerTheme :one
SELECT * FROM server_theme WHERE id = 'singleton';

-- name: UpsertServerTheme :exec
INSERT INTO server_theme (id, colors, fonts, layout, bg_type, generated_css, css_generated_at, updated_by, updated_at)
VALUES ('singleton', $1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id) DO UPDATE SET
  colors = EXCLUDED.colors, fonts = EXCLUDED.fonts, layout = EXCLUDED.layout,
  bg_type = EXCLUDED.bg_type, generated_css = EXCLUDED.generated_css,
  css_generated_at = EXCLUDED.css_generated_at, updated_by = EXCLUDED.updated_by,
  updated_at = EXCLUDED.updated_at;

-- name: CreatePostStyle :one
INSERT INTO post_styles (id, account_id, background_color, background_image_id, text_color, font_family, font_size, font_weight, border_radius, border_color, border_width, padding, has_texture, texture_type, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
RETURNING *;

-- name: GetPostStyleByID :one
SELECT * FROM post_styles WHERE id = $1;

-- name: ListPostStylesByAccount :many
SELECT * FROM post_styles
WHERE account_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: DeletePostStyle :exec
DELETE FROM post_styles WHERE id = $1;
