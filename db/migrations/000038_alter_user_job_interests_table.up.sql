ALTER TABLE "user_job_interests"
DROP COLUMN "job_title";

ALTER TABLE "user_job_interests"
ADD COLUMN "job_position_id" BIGINT;

CREATE INDEX idx_user_job_interests_job_position_id ON "user_job_interests" ("job_position_id");
ALTER TABLE "user_job_interests" ADD FOREIGN KEY ("job_position_id") REFERENCES "job_positions" ("id");