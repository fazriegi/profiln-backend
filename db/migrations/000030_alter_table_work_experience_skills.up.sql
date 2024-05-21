ALTER TABLE "work_experience_skills"
DROP CONSTRAINT work_experience_skills_skill_id_fkey;

ALTER TABLE "work_experience_skills"
RENAME COLUMN "skill_id" TO "user_skill_id";

ALTER TABLE "work_experience_skills"
ADD CONSTRAINT work_experience_skills_user_skill_id_fkey
FOREIGN KEY ("user_skill_id") REFERENCES "user_skills" (id);