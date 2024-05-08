CREATE TABLE "education_skills" (
  "id" BIGSERIAL PRIMARY KEY,
  "education_id" BIGINT,
  "skill_id" BIGINT
);

CREATE INDEX idx_education_skills_id ON "education_skills" ("id");
CREATE INDEX idx_education_skills_education_id ON "education_skills" ("education_id");
CREATE INDEX idx_education_skills_skill_id ON "education_skills" ("skill_id");

ALTER TABLE "education_skills" ADD FOREIGN KEY ("education_id") REFERENCES "educations" ("id");
ALTER TABLE "education_skills" ADD FOREIGN KEY ("skill_id") REFERENCES "skills" ("id");