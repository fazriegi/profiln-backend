CREATE TABLE "user_job_interests" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "job_title" TEXT
);

CREATE INDEX idx_user_job_interests_id ON "user_job_interests" ("id");
CREATE INDEX idx_user_job_interests_user_id ON "user_job_interests" ("user_id");

ALTER TABLE "user_job_interests" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");