-- Drop image_uri column and add title column to about table
ALTER TABLE "about" DROP COLUMN IF EXISTS "image_uri";
ALTER TABLE "about" ADD COLUMN "title" varchar(255);
