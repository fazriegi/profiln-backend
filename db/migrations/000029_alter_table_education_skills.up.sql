ALTER TABLE "education_skills"
DROP CONSTRAINT education_skills_skill_id_fkey;

ALTER TABLE "education_skills"
RENAME COLUMN "skill_id" TO "user_skill_id";

ALTER TABLE "education_skills"
ADD CONSTRAINT education_skills_user_skill_id_fkey
FOREIGN KEY ("user_skill_id") REFERENCES "user_skills" (id);