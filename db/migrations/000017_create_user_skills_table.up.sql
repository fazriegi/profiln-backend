CREATE TABLE "user_skills" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "skill_id" BIGINT,
  "main_skill" BOOL DEFAULT FALSE
);

CREATE INDEX idx_user_skills_id ON "user_skills" ("id");
CREATE INDEX idx_user_skills_user_id ON "user_skills" ("user_id");
CREATE INDEX idx_user_skills_main_skill ON "user_skills" ("main_skill");

ALTER TABLE "user_skills" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "user_skills" ADD FOREIGN KEY ("skill_id") REFERENCES "skills" ("id");