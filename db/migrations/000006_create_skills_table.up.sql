CREATE TABLE "skills" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" TEXT
);

CREATE INDEX idx_skills_id ON "skills" ("id");