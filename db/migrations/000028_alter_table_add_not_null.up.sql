ALTER TABLE "issuing_organizations"
ALTER COLUMN "name" SET NOT NULL;

ALTER TABLE "schools"
ALTER COLUMN "name" SET NOT NULL;

ALTER TABLE "skills"
ALTER COLUMN "name" SET NOT NULL;

ALTER TABLE "companies"
ALTER COLUMN "name" SET NOT NULL;

ALTER TABLE "employment_types"
ALTER COLUMN "name" SET NOT NULL;

ALTER TABLE "location_types"
ALTER COLUMN "name" SET NOT NULL;

ALTER TABLE "social_links"
ALTER COLUMN "name" SET NOT NULL;