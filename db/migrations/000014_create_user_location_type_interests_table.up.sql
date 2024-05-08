CREATE TABLE "user_location_type_interests" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" BIGINT,
  "location_type_id" SMALLINT
);

CREATE INDEX idx_user_location_type_interests_id ON "user_location_type_interests" ("id");
CREATE INDEX idx_user_location_type_interests_user_id ON "user_location_type_interests" ("user_id");

ALTER TABLE "user_location_type_interests" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "user_location_type_interests" ADD FOREIGN KEY ("location_type_id") REFERENCES "location_types" ("id");