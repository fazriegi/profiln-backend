ALTER TABLE "posts"
DROP COLUMN "image_url";

CREATE TABLE "post_images" (
  "id" BIGSERIAL PRIMARY KEY,
  "post_id" BIGINT,
  "url" TEXT
);

CREATE INDEX idx_post_images_id ON "post_images" ("id");
CREATE INDEX idx_post_images_post_id ON "post_images" ("post_id");

ALTER TABLE "post_images" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");