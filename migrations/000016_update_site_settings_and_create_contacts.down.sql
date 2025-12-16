-- Rollback: Drop contacts table
DROP TABLE IF EXISTS "contacts";

-- Rollback: Remove new columns from site_settings
ALTER TABLE "site_settings" DROP COLUMN IF EXISTS "github_url";
ALTER TABLE "site_settings" DROP COLUMN IF EXISTS "youtube_url";
ALTER TABLE "site_settings" DROP COLUMN IF EXISTS "instagram_url";
ALTER TABLE "site_settings" DROP COLUMN IF EXISTS "linkedin_url";
ALTER TABLE "site_settings" DROP COLUMN IF EXISTS "twitter_url";
ALTER TABLE "site_settings" DROP COLUMN IF EXISTS "facebook_url";
ALTER TABLE "site_settings" DROP COLUMN IF EXISTS "logo_big";
ALTER TABLE "site_settings" DROP COLUMN IF EXISTS "site_title";

-- Rollback: Restore old columns
ALTER TABLE "site_settings" ADD COLUMN IF NOT EXISTS "logo_footer" varchar(255);
ALTER TABLE "site_settings" ADD COLUMN IF NOT EXISTS "social_links" jsonb;
ALTER TABLE "site_settings" ADD COLUMN IF NOT EXISTS "contact_info" jsonb;
