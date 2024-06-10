CREATE TABLE "liked_post_comments" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "post_comment_id" BIGINT
);

CREATE INDEX idx_liked_post_comments_id ON "liked_post_comments" ("id");
CREATE INDEX idx_liked_post_comments_user_id ON "liked_post_comments" ("user_id");
CREATE INDEX idx_liked_post_comments_post_comment_id ON "liked_post_comments" ("post_comment_id");

ALTER TABLE "liked_post_comments"
ADD CONSTRAINT liked_post_comments_user_id_post_comment_id_unique UNIQUE ("user_id", "post_comment_id");

ALTER TABLE "liked_post_comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "liked_post_comments" ADD FOREIGN KEY ("post_comment_id") REFERENCES "post_comments" ("id");