ALTER TABLE "educations"
DROP COLUMN "document_url";

CREATE TABLE "education_files" (
  "id" BIGSERIAL PRIMARY KEY,
  "education_id" BIGINT,
  "url" TEXT
);

CREATE INDEX idx_education_files_id ON "education_files" ("id");
CREATE INDEX idx_education_files_education_id ON "education_files" ("education_id");

ALTER TABLE "education_files" ADD FOREIGN KEY ("education_id") REFERENCES "educations" ("id");