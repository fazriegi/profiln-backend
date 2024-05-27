CREATE TABLE "job_positions" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" TEXT UNIQUE
);

CREATE INDEX idx_job_positions_id ON "job_positions" ("id");
CREATE INDEX idx_job_positions_name ON "job_positions" ("name");