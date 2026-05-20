-- name: GetAccountByID :one
SELECT * FROM accounts WHERE id = $1;

-- name: GetAccountByURI :one
SELECT * FROM accounts WHERE uri = $1;

-- name: GetAccountByUsername :one
SELECT * FROM accounts WHERE username = $1 AND domain IS NULL AND is_local = TRUE;

-- name: GetAccountByUsernameAndDomain :one
SELECT * FROM accounts WHERE username = $1 AND domain = $2;

-- name: GetAccountByEmail :one
SELECT * FROM accounts WHERE email = $1 AND is_local = TRUE;

-- name: GetAccountByCustomDomain :one
SELECT * FROM accounts WHERE custom_domain = $1 AND custom_domain_verified = TRUE;

-- name: GetAccountByEmailVerifyToken :one
SELECT * FROM accounts WHERE email_verify_token = $1 AND is_local = TRUE;

-- name: CreateAccount :one
INSERT INTO accounts (
    id, username, domain, uri, display_name, bio, bio_text,
    public_key_pem, private_key_pem, key_id, role,
    is_local, actor_type,
    inbox_url, outbox_url, shared_inbox_url,
    followers_url, following_url, featured_url,
    password_hash, email, email_verify_token
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11,
    $12, $13,
    $14, $15, $16,
    $17, $18, $19,
    $20, $21, $22
) RETURNING *;

-- name: UpdateAccountProfile :one
UPDATE accounts SET
    display_name = COALESCE($2, display_name),
    bio = COALESCE($3, bio),
    bio_text = COALESCE($4, bio_text),
    avatar_media_id = COALESCE($5, avatar_media_id),
    header_media_id = COALESCE($6, header_media_id),
    is_locked = COALESCE($7, is_locked),
    is_bot = COALESCE($8, is_bot),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateAccountSuspended :exec
UPDATE accounts SET is_suspended = $2, updated_at = NOW() WHERE id = $1;

-- name: UpdateAccountSilenced :exec
UPDATE accounts SET is_silenced = $2, updated_at = NOW() WHERE id = $1;

-- name: UpdateAccountRole :exec
UPDATE accounts SET role = $2, updated_at = NOW() WHERE id = $1;

-- name: VerifyAccountEmail :exec
UPDATE accounts SET email_verified_at = NOW(), email_verify_token = NULL, updated_at = NOW() WHERE id = $1;

-- name: UpdateAccountPassword :exec
UPDATE accounts SET password_hash = $2, updated_at = NOW() WHERE id = $1;

-- name: IncrementFollowersCount :exec
UPDATE accounts SET followers_count = followers_count + 1 WHERE id = $1;

-- name: DecrementFollowersCount :exec
UPDATE accounts SET followers_count = GREATEST(followers_count - 1, 0) WHERE id = $1;

-- name: IncrementFollowingCount :exec
UPDATE accounts SET following_count = following_count + 1 WHERE id = $1;

-- name: DecrementFollowingCount :exec
UPDATE accounts SET following_count = GREATEST(following_count - 1, 0) WHERE id = $1;

-- name: IncrementPostsCount :exec
UPDATE accounts SET posts_count = posts_count + 1 WHERE id = $1;

-- name: DecrementPostsCount :exec
UPDATE accounts SET posts_count = GREATEST(posts_count - 1, 0) WHERE id = $1;

-- name: ListLocalAccounts :many
SELECT * FROM accounts WHERE is_local = TRUE ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: SearchAccountsByUsername :many
SELECT * FROM accounts
WHERE username ILIKE $1 || '%'
ORDER BY is_local DESC, followers_count DESC
LIMIT $2;
