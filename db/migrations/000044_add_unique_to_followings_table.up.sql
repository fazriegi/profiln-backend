ALTER TABLE "followings"
ADD CONSTRAINT followings_user_id_follow_user_id_unique UNIQUE ("user_id", "follow_user_id");