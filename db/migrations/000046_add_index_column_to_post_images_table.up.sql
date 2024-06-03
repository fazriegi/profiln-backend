ALTER TABLE "post_images"
ADD COLUMN "index" SMALLINT DEFAULT 0;

CREATE INDEX idx_post_images_index ON "post_images" ("index");