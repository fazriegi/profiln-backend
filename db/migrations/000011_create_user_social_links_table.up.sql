CREATE TABLE "user_social_links" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "social_link_id" SMALLINT,
  "url" TEXT
);

CREATE INDEX idx_user_social_links_id ON "user_social_links" ("id");
CREATE INDEX idx_user_social_links_user_id ON "user_social_links" ("user_id");

ALTER TABLE "user_social_links" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "user_social_links" ADD FOREIGN KEY ("social_link_id") REFERENCES "user_social_links" ("id");