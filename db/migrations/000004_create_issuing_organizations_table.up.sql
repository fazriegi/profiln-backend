CREATE TABLE "issuing_organizations" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" TEXT
);

CREATE INDEX idx_issuing_organizations_id ON "issuing_organizations" ("id");
CREATE INDEX idx_issuing_organizations_name ON "issuing_organizations" ("name");