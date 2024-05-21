ALTER TABLE "education_skills"
DROP CONSTRAINT education_skills_user_skill_id_fkey;

ALTER TABLE "education_skills"
RENAME COLUMN "user_skill_id" TO "skill_id";

ALTER TABLE "education_skills"
ADD CONSTRAINT education_skills_skill_id_fkey
FOREIGN KEY ("skill_id") REFERENCES "skills" (id);
