-- DROP TYPE IF EXISTS "subscriber_status";
-- DROP TYPE IF EXISTS "comment_status";
-- DROP TYPE IF EXISTS "post_status";

DROP TYPE IF EXISTS "post_status";
DROP TYPE IF EXISTS "comment_status";
DROP TYPE IF EXISTS "subscriber_status";

CREATE TYPE "post_status" AS ENUM ('draft', 'published', 'archived');
CREATE TYPE "comment_status" AS ENUM ('pending', 'approved', 'spam');
CREATE TYPE "subscriber_status" AS ENUM ('active', 'unsubscribed', 'bounced');