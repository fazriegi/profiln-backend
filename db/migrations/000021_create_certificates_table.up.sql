CREATE TABLE "certificates" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "name" TEXT,
  "issuing_organization_id" BIGINT,
  "issue_date" DATE,
  "expiration_date" DATE,
  "credential_id" TEXT,
  "url" TEXT,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_certificates_id ON "certificates" ("id");
CREATE INDEX idx_certificates_user_id ON "certificates" ("user_id");
CREATE INDEX idx_certificates_issuing_organization_id ON "certificates" ("issuing_organization_id");

ALTER TABLE "certificates" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "certificates" ADD FOREIGN KEY ("issuing_organization_id") REFERENCES "issuing_organizations" ("id");