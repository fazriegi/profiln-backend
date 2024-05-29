ALTER TABLE "posts"
ADD COLUMN "repost" BOOL DEFAULT FALSE,
ADD COLUMN "original_post_id" BIGINT;

ALTER TABLE "posts"
DROP COLUMN "title",
DROP COLUMN "visibility";

DROP INDEX IF EXISTS idx_posts_title;
DROP INDEX IF EXISTS idx_posts_visibility;

ALTER TABLE "posts" ADD FOREIGN KEY ("original_post_id") REFERENCES "posts" ("id");