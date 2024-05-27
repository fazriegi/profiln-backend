DROP INDEX IF EXISTS idx_user_job_interests;

ALTER TABLE "user_job_interests"
DROP COLUMN "job_position_id";

ALTER TABLE "user_job_interests"
ADD COLUMN "job_title" TEXT;