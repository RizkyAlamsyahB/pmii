-- Rollback: Remove title and restore image_uri
ALTER TABLE "about" DROP COLUMN IF EXISTS "title";
ALTER TABLE "about" ADD COLUMN "image_uri" varchar(255);
