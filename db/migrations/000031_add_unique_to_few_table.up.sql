ALTER TABLE "user_social_links"
ADD CONSTRAINT user_social_links_user_id_social_link_id_unique UNIQUE ("user_id", "social_link_id");

ALTER TABLE "user_skills"
ADD CONSTRAINT user_skills_user_id_skill_id_unique UNIQUE ("user_id", "skill_id");

ALTER TABLE "work_experience_skills"
ADD CONSTRAINT work_experience_skills_work_experience_id_user_skill_id_unique UNIQUE ("work_experience_id", "user_skill_id");

ALTER TABLE "education_skills"
ADD CONSTRAINT education_skills_education_id_user_skill_id_unique UNIQUE ("education_id", "user_skill_id");