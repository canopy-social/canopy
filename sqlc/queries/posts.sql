-- name: CreatePost :one
INSERT INTO posts (
    id, uri, url, account_id, content, content_text,
    content_warning, is_sensitive, visibility, language,
    reply_to_id, reply_to_uri, thread_root_id,
    boost_of_id, boost_of_uri, is_local, post_style_id
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10,
    $11, $12, $13,
    $14, $15, $16, $17
) RETURNING *;

-- name: GetPostByID :one
SELECT * FROM posts WHERE id = $1;

-- name: GetPostByURI :one
SELECT * FROM posts WHERE uri = $1;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;

-- name: UpdatePostContent :one
UPDATE posts SET
    content = $2,
    content_text = $3,
    content_warning = $4,
    is_sensitive = $5,
    edit_history = edit_history || $6::jsonb,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: PinPost :exec
UPDATE posts SET is_pinned = TRUE WHERE id = $1 AND account_id = $2;

-- name: UnpinPost :exec
UPDATE posts SET is_pinned = FALSE WHERE id = $1 AND account_id = $2;

-- name: ListPostsByAccount :many
SELECT * FROM posts
WHERE account_id = $1 AND visibility != 'direct' AND boost_of_id IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPostsByAccountWithBoosts :many
SELECT * FROM posts
WHERE account_id = $1 AND visibility != 'direct'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPinnedPosts :many
SELECT * FROM posts
WHERE account_id = $1 AND is_pinned = TRUE
ORDER BY created_at DESC;

-- name: ListPublicTimeline :many
SELECT * FROM posts
WHERE visibility = 'public' AND is_local = TRUE AND reply_to_id IS NULL AND boost_of_id IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListThreadReplies :many
SELECT * FROM posts
WHERE reply_to_id = $1
ORDER BY created_at ASC;

-- name: ListThreadContext :many
SELECT * FROM posts
WHERE thread_root_id = $1
ORDER BY created_at ASC;

-- name: IncrementLikesCount :exec
UPDATE posts SET likes_count = likes_count + 1 WHERE id = $1;

-- name: DecrementLikesCount :exec
UPDATE posts SET likes_count = GREATEST(likes_count - 1, 0) WHERE id = $1;

-- name: IncrementBoostsCount :exec
UPDATE posts SET boosts_count = boosts_count + 1 WHERE id = $1;

-- name: DecrementBoostsCount :exec
UPDATE posts SET boosts_count = GREATEST(boosts_count - 1, 0) WHERE id = $1;

-- name: IncrementRepliesCount :exec
UPDATE posts SET replies_count = replies_count + 1 WHERE id = $1;

-- name: DecrementRepliesCount :exec
UPDATE posts SET replies_count = GREATEST(replies_count - 1, 0) WHERE id = $1;

-- name: CreatePostLike :exec
INSERT INTO post_likes (id, post_id, account_id, uri)
VALUES ($1, $2, $3, $4)
ON CONFLICT (post_id, account_id) DO NOTHING;

-- name: DeletePostLike :exec
DELETE FROM post_likes WHERE post_id = $1 AND account_id = $2;

-- name: HasLikedPost :one
SELECT EXISTS(SELECT 1 FROM post_likes WHERE post_id = $1 AND account_id = $2);

-- name: CreatePostBoost :exec
INSERT INTO post_boosts (id, post_id, account_id, uri)
VALUES ($1, $2, $3, $4)
ON CONFLICT (post_id, account_id) DO NOTHING;

-- name: DeletePostBoost :exec
DELETE FROM post_boosts WHERE post_id = $1 AND account_id = $2;

-- name: HasBoostedPost :one
SELECT EXISTS(SELECT 1 FROM post_boosts WHERE post_id = $1 AND account_id = $2);

-- name: CreatePostMention :exec
INSERT INTO post_mentions (post_id, account_id, uri)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: ListPostMentions :many
SELECT a.* FROM accounts a
INNER JOIN post_mentions pm ON pm.account_id = a.id
WHERE pm.post_id = $1;

-- name: CreatePostTag :exec
INSERT INTO post_tags (post_id, tag)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: ListPostTags :many
SELECT tag FROM post_tags WHERE post_id = $1;

-- name: SearchPostsByTag :many
SELECT p.* FROM posts p
INNER JOIN post_tags pt ON pt.post_id = p.id
WHERE pt.tag = $1 AND p.visibility = 'public'
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;
