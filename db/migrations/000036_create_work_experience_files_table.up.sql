ALTER TABLE "work_experiences"
DROP COLUMN "image_url";

CREATE TABLE "work_experience_files" (
  "id" BIGSERIAL PRIMARY KEY,
  "work_experience_id" BIGINT,
  "url" TEXT
);

CREATE INDEX idx_work_experience_files_id ON "work_experience_files" ("id");
CREATE INDEX idx_work_experience_files_work_experience_id ON "work_experience_files" ("work_experience_id");

ALTER TABLE "work_experience_files" ADD FOREIGN KEY ("work_experience_id") REFERENCES "work_experiences" ("id");