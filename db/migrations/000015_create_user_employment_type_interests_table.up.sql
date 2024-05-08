CREATE TABLE "user_employment_type_interests" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "employment_type_id" SMALLINT
);

CREATE INDEX idx_user_employment_type_interests_id ON "user_employment_type_interests" ("id");
CREATE INDEX idx_user_employment_type_interests_user_id ON "user_employment_type_interests" ("user_id");

ALTER TABLE "user_employment_type_interests" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "user_employment_type_interests" ADD FOREIGN KEY ("employment_type_id") REFERENCES "employment_types" ("id");