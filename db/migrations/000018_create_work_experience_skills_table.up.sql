CREATE TABLE "work_experience_skills" (
  "id" BIGSERIAL PRIMARY KEY,
  "work_experience_id" BIGINT,
  "skill_id" BIGINT
);

CREATE INDEX idx_work_experience_skills_id ON "work_experience_skills" ("id");
CREATE INDEX idx_work_experience_skills_work_experience_id ON "work_experience_skills" ("work_experience_id");
CREATE INDEX idx_work_experience_skills_skill_id ON "work_experience_skills" ("skill_id");

ALTER TABLE "work_experience_skills" ADD FOREIGN KEY ("work_experience_id") REFERENCES "work_experiences" ("id");
ALTER TABLE "work_experience_skills" ADD FOREIGN KEY ("skill_id") REFERENCES "skills" ("id");