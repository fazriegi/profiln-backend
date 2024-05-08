CREATE TABLE "companies" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" TEXT
);

CREATE INDEX idx_companies_id ON "companies" ("id");
CREATE INDEX idx_companies_name ON "companies" ("name");