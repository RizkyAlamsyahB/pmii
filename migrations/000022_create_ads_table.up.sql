-- Create enum for ad pages
CREATE TYPE "ad_page" AS ENUM (
  'landing',
  'news',
  'opini',
  'life_at_pmii',
  'islamic',
  'detail_article'
);

-- Create ads table
CREATE TABLE IF NOT EXISTS "ads" (
  "id" SERIAL PRIMARY KEY,
  "page" ad_page NOT NULL,
  "slot" INTEGER NOT NULL,
  "image_url" VARCHAR(500),
  "resolution" VARCHAR(20) NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  
  CONSTRAINT "unique_page_slot" UNIQUE ("page", "slot")
);

-- Create index for faster queries
CREATE INDEX "idx_ads_page" ON "ads" ("page");

-- Insert default ad slots based on design
-- Landing Page
INSERT INTO "ads" ("page", "slot", "resolution") VALUES 
  ('landing', 1, '728x90'),
  ('landing', 2, '16x9');

-- News Page
INSERT INTO "ads" ("page", "slot", "resolution") VALUES 
  ('news', 1, '16x9'),
  ('news', 2, '4x3');

-- Opini Page
INSERT INTO "ads" ("page", "slot", "resolution") VALUES 
  ('opini', 1, '3x4'),
  ('opini', 2, '16x9'),
  ('opini', 3, '4x3');

-- Life at PMII Page
INSERT INTO "ads" ("page", "slot", "resolution") VALUES 
  ('life_at_pmii', 1, '16x9'),
  ('life_at_pmii', 2, '4x3');

-- Islamic Page
INSERT INTO "ads" ("page", "slot", "resolution") VALUES 
  ('islamic', 1, '3x4'),
  ('islamic', 2, '16x9'),
  ('islamic', 3, '4x3');

-- Detail Article Page
INSERT INTO "ads" ("page", "slot", "resolution") VALUES 
  ('detail_article', 1, '3x4'),
  ('detail_article', 2, '9x16'),
  ('detail_article', 3, '728x90');
