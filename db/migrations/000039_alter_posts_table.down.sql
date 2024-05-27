ALTER TABLE "posts"
DROP COLUMN "title",
DROP COLUMN "visibility";

DROP INDEX IF EXISTS idx_posts_title;
DROP INDEX IF EXISTS idx_posts_visibility;