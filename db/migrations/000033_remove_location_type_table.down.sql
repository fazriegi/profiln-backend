CREATE TABLE "location_types" (
  "id" SMALLSERIAL PRIMARY KEY,
  "name" VARCHAR(10)
);

-- add location type id to user location type interests table
ALTER TABLE "user_location_type_interests"
DROP COLUMN "location_type";

ALTER TABLE "user_location_type_interests"
ADD COLUMN "location_type_id" BIGINT;

DROP INDEX IF EXISTS idx_user_location_type_interests_location_type;
CREATE INDEX idx_user_location_type_interests_user_location_type_id ON "user_location_type_interests" ("location_type_id");
ALTER TABLE "user_location_type_interests" ADD FOREIGN KEY ("location_type_id") REFERENCES "location_types" ("id");


-- add location type id to work experiences table
ALTER TABLE "work_experiences"
DROP COLUMN "location_type";
ALTER TABLE "work_experiences"
ADD COLUMN "location_type_id" BIGINT;
DROP INDEX IF EXISTS idx_work_experiences_location_type;
CREATE INDEX idx_work_experiences_location_type_id ON "work_experiences" ("location_type_id");
ALTER TABLE "work_experiences" ADD FOREIGN KEY ("location_type_id") REFERENCES "location_types" ("id");