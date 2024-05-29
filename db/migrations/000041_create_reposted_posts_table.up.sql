CREATE TABLE "reposted_posts" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "post_id" BIGINT
);

CREATE INDEX idx_reposted_posts_id ON "reposted_posts" ("id");
CREATE INDEX idx_reposted_posts_user_id ON "reposted_posts" ("user_id");
CREATE INDEX idx_reposted_posts_post_id ON "reposted_posts" ("post_id");

ALTER TABLE "reposted_posts"
ADD CONSTRAINT reposted_posts_user_id_post_id_unique UNIQUE ("user_id", "post_id");

ALTER TABLE "reposted_posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "reposted_posts" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");