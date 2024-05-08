CREATE TABLE "educations" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "school_id" BIGINT,
  "degree" VARCHAR(50),
  "field_of_study" VARCHAR(255),
  "gpa" NUMERIC(3, 2),
  "start_date" DATE,
  "finish_date" DATE,
  "description" TEXT,
  "document_url" TEXT,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_educations_id ON "educations" ("id");
CREATE INDEX idx_educations_user_id ON "educations" ("user_id");
CREATE INDEX idx_educations_school_id ON "educations" ("school_id");

ALTER TABLE "educations" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "educations" ADD FOREIGN KEY ("school_id") REFERENCES "schools" ("id");