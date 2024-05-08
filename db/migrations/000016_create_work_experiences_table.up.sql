CREATE TABLE "work_experiences" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "job_title" TEXT,
  "company_id" BIGINT,
  "employment_type_id" SMALLINT,
  "location" TEXT,
  "location_type_id" SMALLINT,
  "start_date" DATE,
  "finish_date" DATE,
  "description" TEXT,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_work_experiences_id ON "work_experiences" ("id");
CREATE INDEX idx_work_experiences_user_id ON "work_experiences" ("user_id");
CREATE INDEX idx_work_experiences_job_title ON "work_experiences" ("job_title");
CREATE INDEX idx_work_experiences_company_id ON "work_experiences" ("company_id");

ALTER TABLE "work_experiences" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "work_experiences" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");
ALTER TABLE "work_experiences" ADD FOREIGN KEY ("employment_type_id") REFERENCES "employment_types" ("id");
ALTER TABLE "work_experiences" ADD FOREIGN KEY ("location_type_id") REFERENCES "location_types" ("id");