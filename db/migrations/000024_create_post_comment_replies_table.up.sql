CREATE TABLE "post_comment_replies" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "post_comment_id" BIGINT,
  "content" TEXT,
  "image_url" TEXT,
  "like_count" INTEGER DEFAULT 0,
  "is_post_author" BOOL DEFAULT FALSE,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_post_comment_replies_id ON "post_comment_replies" ("id");
CREATE INDEX idx_post_comment_replies_user_id ON "post_comment_replies" ("user_id");
CREATE INDEX idx_post_comment_replies_post_comment_id ON "post_comment_replies" ("post_comment_id");

ALTER TABLE "post_comment_replies" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "post_comment_replies" ADD FOREIGN KEY ("post_comment_id") REFERENCES "post_comments" ("id");