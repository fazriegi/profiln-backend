CREATE TABLE "user_details" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "phone_number" VARCHAR(20),
  "gender" VARCHAR(1),
  "location" TEXT,
  "portfolio_url" TEXT,
  "about" TEXT,
  "hide_phone_number" BOOL DEFAULT FALSE,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_details_id ON "user_details" ("id");
CREATE INDEX idx_user_details_user_id ON "user_details" ("user_id");

ALTER TABLE "user_details" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");