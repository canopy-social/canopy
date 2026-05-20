-- 0005_create_essays.down.sql
ALTER TABLE media_attachments DROP CONSTRAINT IF EXISTS fk_media_essay;
DROP TABLE IF EXISTS essay_marginalia;
DROP TABLE IF EXISTS essays;
