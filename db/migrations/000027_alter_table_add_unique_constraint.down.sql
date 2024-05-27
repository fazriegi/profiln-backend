ALTER TABLE "issuing_organizations"
DROP CONSTRAINT unique_issuing_organizations_name;

ALTER TABLE "schools"
DROP CONSTRAINT unique_schools_name;

ALTER TABLE "skills"
DROP CONSTRAINT unique_skills_name;

ALTER TABLE "companies"
DROP CONSTRAINT unique_companies_name;

ALTER TABLE "employment_types"
DROP CONSTRAINT unique_employment_types_name;

ALTER TABLE "location_types"
DROP CONSTRAINT unique_location_types_name;

ALTER TABLE "social_links"
DROP CONSTRAINT unique_social_links_name;