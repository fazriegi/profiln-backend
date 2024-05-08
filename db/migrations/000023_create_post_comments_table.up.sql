CREATE TABLE "post_comments" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "post_id" BIGINT,
  "content" TEXT,
  "image_url" TEXT,
  "like_count" INTEGER,
  "reply_count" INTEGER,
  "is_post_author" BOOL DEFAULT FALSE,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_post_comments_id ON "post_comments" ("id");
CREATE INDEX idx_post_comments_user_id ON "post_comments" ("user_id");
CREATE INDEX idx_post_comments_post_id ON "post_comments" ("post_id");

ALTER TABLE "post_comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "post_comments" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");