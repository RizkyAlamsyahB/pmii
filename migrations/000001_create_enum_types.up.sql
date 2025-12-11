CREATE TYPE "post_status" AS ENUM (
  'draft',
  'published',
  'archived'
);

CREATE TYPE "comment_status" AS ENUM (
  'pending',
  'approved',
  'spam'
);

CREATE TYPE "subscriber_status" AS ENUM (
  'active',
  'unsubscribed',
  'bounced'
);
