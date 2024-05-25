-- remove location type id from user location type interests table
ALTER TABLE "user_location_type_interests"
DROP COLUMN "location_type_id";

ALTER TABLE "user_location_type_interests"
ADD COLUMN "location_type" VARCHAR(10);

DROP INDEX IF EXISTS idx_user_location_type_interests_user_id;
CREATE INDEX idx_user_location_type_interests_location_type ON "user_location_type_interests" ("location_type");


-- remove location type id from work experiences table
ALTER TABLE "work_experiences"
DROP COLUMN "location_type_id";

ALTER TABLE "work_experiences"
ADD COLUMN "location_type" VARCHAR(10);

CREATE INDEX idx_work_experiences_location_type ON "work_experiences" ("location_type");

-- remove location types table
DROP TABLE "location_types";