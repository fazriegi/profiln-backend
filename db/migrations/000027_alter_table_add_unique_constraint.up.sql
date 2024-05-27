ALTER TABLE "issuing_organizations"
ADD CONSTRAINT unique_issuing_organizations_name UNIQUE ("name");

ALTER TABLE "schools"
ADD CONSTRAINT unique_schools_name UNIQUE ("name");

ALTER TABLE "skills"
ADD CONSTRAINT unique_skills_name UNIQUE ("name");

ALTER TABLE "companies"
ADD CONSTRAINT unique_companies_name UNIQUE ("name");

ALTER TABLE "employment_types"
ADD CONSTRAINT unique_employment_types_name UNIQUE ("name");

ALTER TABLE "location_types"
ADD CONSTRAINT unique_location_types_name UNIQUE ("name");

ALTER TABLE "social_links"
ADD CONSTRAINT unique_social_links_name UNIQUE ("name");