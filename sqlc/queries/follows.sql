-- name: CreateFollow :one
INSERT INTO follows (id, follower_id, following_id, status, uri)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFollow :one
SELECT * FROM follows WHERE follower_id = $1 AND following_id = $2;

-- name: GetFollowByURI :one
SELECT * FROM follows WHERE uri = $1;

-- name: UpdateFollowStatus :exec
UPDATE follows SET status = $2 WHERE id = $1;

-- name: DeleteFollow :exec
DELETE FROM follows WHERE follower_id = $1 AND following_id = $2;

-- name: DeleteFollowByID :exec
DELETE FROM follows WHERE id = $1;

-- name: ListFollowers :many
SELECT a.* FROM accounts a
INNER JOIN follows f ON f.follower_id = a.id
WHERE f.following_id = $1 AND f.status = 'accepted'
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListFollowing :many
SELECT a.* FROM accounts a
INNER JOIN follows f ON f.following_id = a.id
WHERE f.follower_id = $1 AND f.status = 'accepted'
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPendingFollowRequests :many
SELECT a.* FROM accounts a
INNER JOIN follows f ON f.follower_id = a.id
WHERE f.following_id = $1 AND f.status = 'pending'
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountFollowers :one
SELECT COUNT(*) FROM follows WHERE following_id = $1 AND status = 'accepted';

-- name: CountFollowing :one
SELECT COUNT(*) FROM follows WHERE follower_id = $1 AND status = 'accepted';

-- name: IsFollowing :one
SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2 AND status = 'accepted');

-- name: CreateBlock :one
INSERT INTO blocks (id, account_id, target_id, uri)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteBlock :exec
DELETE FROM blocks WHERE account_id = $1 AND target_id = $2;

-- name: IsBlocking :one
SELECT EXISTS(SELECT 1 FROM blocks WHERE account_id = $1 AND target_id = $2);

-- name: ListBlocks :many
SELECT a.* FROM accounts a
INNER JOIN blocks b ON b.target_id = a.id
WHERE b.account_id = $1
ORDER BY b.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateMute :one
INSERT INTO mutes (id, account_id, target_id, hide_notifications)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteMute :exec
DELETE FROM mutes WHERE account_id = $1 AND target_id = $2;

-- name: IsMuting :one
SELECT EXISTS(SELECT 1 FROM mutes WHERE account_id = $1 AND target_id = $2);

-- name: ListMutes :many
SELECT a.* FROM accounts a
INNER JOIN mutes m ON m.target_id = a.id
WHERE m.account_id = $1
ORDER BY m.created_at DESC
LIMIT $2 OFFSET $3;
