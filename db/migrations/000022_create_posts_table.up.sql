CREATE TABLE "posts" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "content" TEXT,
  "image_url" TEXT,
  "like_count" INTEGER,
  "comment_count" INTEGER,
  "repost_count" INTEGER,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_posts_id ON "posts" ("id");
CREATE INDEX idx_posts_user_id ON "posts" ("user_id");

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");