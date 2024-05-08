CREATE TABLE "schools" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" TEXT
);

CREATE INDEX idx_schools_id ON "schools" ("id");
CREATE INDEX idx_schools_name ON "schools" ("name");