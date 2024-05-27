ALTER TABLE "posts"
ADD COLUMN "title" TEXT NOT NULL DEFAULT '',
ADD COLUMN "visibility" VARCHAR(10) NOT NULL DEFAULT 'public';

CREATE INDEX idx_posts_title ON "posts" ("title");
CREATE INDEX idx_posts_visibility ON "posts" ("visibility");