ALTER TABLE "user_social_links"
ADD CONSTRAINT user_social_links_user_id_social_link_id_unique UNIQUE ("user_id", "social_link_id");
