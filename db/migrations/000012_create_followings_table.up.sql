CREATE TABLE "followings" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "follow_user_id" BIGINT
);

CREATE INDEX idx_followings_user_id ON "followings" ("user_id");
CREATE INDEX idx_followings_follow_user_id ON "followings" ("follow_user_id");

ALTER TABLE "followings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "followings" ADD FOREIGN KEY ("follow_user_id") REFERENCES "users" ("id");