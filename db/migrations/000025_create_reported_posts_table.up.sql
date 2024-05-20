CREATE TABLE "reported_posts" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "post_id" BIGINT,
  "reason" VARCHAR(15),
  "message" TEXT
);

CREATE INDEX idx_reported_posts_id ON "reported_posts" ("id");
CREATE INDEX idx_reported_posts_user_id ON "reported_posts" ("user_id");
CREATE INDEX idx_reported_posts_post_id ON "reported_posts" ("post_id");
CREATE INDEX idx_reported_posts_reason ON "reported_posts" ("reason");

ALTER TABLE "reported_posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "reported_posts" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");