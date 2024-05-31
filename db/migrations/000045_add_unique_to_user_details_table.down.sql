ALTER TABLE "user_details"
DROP CONSTRAINT user_details_user_id_unique;

ALTER TABLE "user_details"
ALTER COLUMN "user_id" DROP NOT NULL;