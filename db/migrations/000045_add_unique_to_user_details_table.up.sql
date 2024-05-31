ALTER TABLE "user_details"
ADD CONSTRAINT user_details_user_id_unique UNIQUE ("user_id");

ALTER TABLE "user_details"
ALTER COLUMN "user_id" SET NOT NULL;