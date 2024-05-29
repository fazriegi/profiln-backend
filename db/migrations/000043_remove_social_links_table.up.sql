ALTER TABLE "user_social_links"
DROP COLUMN "social_link_id";

ALTER TABLE "user_social_links"
ADD COLUMN "platform" VARCHAR(20);

ALTER TABLE "user_social_links"
ADD CONSTRAINT user_social_links_user_id_platform_unique UNIQUE ("user_id", "platform");

CREATE INDEX idx_user_social_links_name ON "user_social_links" ("platform");

DROP TABLE "social_links";