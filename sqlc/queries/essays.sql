-- name: CreateEssay :one
INSERT INTO essays (
    id, uri, url, account_id, title, slug, subtitle,
    content, content_text, content_raw, cover_media_id,
    reading_time_minutes, visibility, language, is_local,
    word_count, published_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11,
    $12, $13, $14, $15,
    $16, $17
) RETURNING *;

-- name: GetEssayByID :one
SELECT * FROM essays WHERE id = $1;

-- name: GetEssayByURI :one
SELECT * FROM essays WHERE uri = $1;

-- name: GetEssayBySlug :one
SELECT * FROM essays WHERE account_id = $1 AND slug = $2;

-- name: UpdateEssay :one
UPDATE essays SET
    title = COALESCE($2, title),
    subtitle = COALESCE($3, subtitle),
    content = COALESCE($4, content),
    content_text = COALESCE($5, content_text),
    content_raw = COALESCE($6, content_raw),
    cover_media_id = COALESCE($7, cover_media_id),
    reading_time_minutes = COALESCE($8, reading_time_minutes),
    visibility = COALESCE($9, visibility),
    language = COALESCE($10, language),
    word_count = COALESCE($11, word_count),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: PublishEssay :one
UPDATE essays SET published_at = NOW(), ap_published = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UnpublishEssay :one
UPDATE essays SET published_at = NULL, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteEssay :exec
DELETE FROM essays WHERE id = $1;

-- name: ListEssaysByAccount :many
SELECT * FROM essays
WHERE account_id = $1 AND published_at IS NOT NULL
ORDER BY published_at DESC
LIMIT $2 OFFSET $3;

-- name: ListDraftsByAccount :many
SELECT * FROM essays
WHERE account_id = $1 AND published_at IS NULL
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: IncrementEssayViews :exec
UPDATE essays SET views_count = views_count + 1 WHERE id = $1;

-- name: IncrementEssayLikes :exec
UPDATE essays SET likes_count = likes_count + 1 WHERE id = $1;

-- name: DecrementEssayLikes :exec
UPDATE essays SET likes_count = GREATEST(likes_count - 1, 0) WHERE id = $1;
