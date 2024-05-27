CREATE TABLE "liked_posts" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "post_id" BIGINT
);

CREATE INDEX idx_liked_posts_id ON "liked_posts" ("id");
CREATE INDEX idx_liked_posts_user_id ON "liked_posts" ("user_id");
CREATE INDEX idx_liked_posts_post_id ON "liked_posts" ("post_id");

ALTER TABLE "liked_posts"
ADD CONSTRAINT liked_posts_user_id_post_id_unique UNIQUE ("user_id", "post_id");

ALTER TABLE "liked_posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "liked_posts" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");