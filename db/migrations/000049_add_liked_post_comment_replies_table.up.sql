CREATE TABLE "liked_post_comment_replies" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "post_comment_reply_id" BIGINT
);

CREATE INDEX idx_liked_post_comment_replies_id ON "liked_post_comment_replies" ("id");
CREATE INDEX idx_liked_post_comment_replies_user_id ON "liked_post_comment_replies" ("user_id");
CREATE INDEX idx_liked_post_comment_replies_post_comment_reply_id ON "liked_post_comment_replies" ("post_comment_reply_id");

ALTER TABLE "liked_post_comment_replies"
ADD CONSTRAINT liked_post_comment_replies_user_id_post_comment_replies_id_unique UNIQUE ("user_id", "post_comment_reply_id");

ALTER TABLE "liked_post_comment_replies" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "liked_post_comment_replies" ADD FOREIGN KEY ("post_comment_reply_id") REFERENCES "post_comment_replies" ("id");