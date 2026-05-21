ALTER TABLE media_attachments DROP CONSTRAINT IF EXISTS fk_media_post;
DROP TABLE IF EXISTS post_boosts;
DROP TABLE IF EXISTS post_likes;
DROP TABLE IF EXISTS post_tags;
DROP TABLE IF EXISTS post_mentions;
DROP TABLE IF EXISTS posts;
