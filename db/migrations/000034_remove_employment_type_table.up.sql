-- remove employment type id from user employment type interests table
ALTER TABLE "user_employment_type_interests"
DROP COLUMN "employment_type_id";

ALTER TABLE "user_employment_type_interests"
ADD COLUMN "employment_type" VARCHAR(20);

CREATE INDEX idx_user_employment_type_interests_employment_type ON "user_employment_type_interests" ("employment_type");


-- remove employment type id from work experiences table
ALTER TABLE "work_experiences"
DROP COLUMN "employment_type_id";

ALTER TABLE "work_experiences"
ADD COLUMN "employment_type" VARCHAR(20);

CREATE INDEX idx_work_experiences_employment_type ON "work_experiences" ("employment_type");

-- remove employment types table
DROP TABLE "employment_types";