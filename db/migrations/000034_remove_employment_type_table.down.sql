CREATE TABLE "employment_types" (
  "id" SMALLSERIAL PRIMARY KEY,
  "name" VARCHAR(10)
);

-- add employment type id to user employment type interests table
ALTER TABLE "user_employment_type_interests"
DROP COLUMN "employment_type";

ALTER TABLE "user_employment_type_interests"
ADD COLUMN "employment_type_id" BIGINT;

DROP INDEX IF EXISTS idx_user_employment_type_interests_employment_type;
ALTER TABLE "user_employment_type_interests" ADD FOREIGN KEY ("employment_type_id") REFERENCES "employment_types" ("id");


-- add employment type id to work experiences table
ALTER TABLE "work_experiences"
DROP COLUMN "employment_type";
ALTER TABLE "work_experiences"
ADD COLUMN "employment_type_id" BIGINT;
DROP INDEX IF EXISTS idx_work_experiences_employment_type;
ALTER TABLE "work_experiences" ADD FOREIGN KEY ("employment_type_id") REFERENCES "employment_types" ("id");