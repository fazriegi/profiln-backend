CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "password" text,
  "full_name" varchar NOT NULL,
  "verified_email" bool DEFAULT false
);

CREATE TABLE "user_otps" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "otp" varchar(6)
);

CREATE INDEX ON "users" ("email");

ALTER TABLE "user_otps" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
