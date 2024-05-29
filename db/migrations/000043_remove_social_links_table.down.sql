ALTER TABLE "user_social_links"
DROP COLUMN "platform";


CREATE TABLE "social_links" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR(20)
);

CREATE INDEX idx_social_links_id ON "social_links" ("id");
CREATE INDEX idx_social_links_name ON "social_links" ("name");

ALTER TABLE "user_social_links"
ADD COLUMN "social_link_id" SMALLINT;

CREATE INDEX idx_user_social_links_social_link_id ON "user_social_links" ("social_link_id");
ALTER TABLE "user_social_links" ADD FOREIGN KEY ("social_link_id") REFERENCES "social_links" ("id");