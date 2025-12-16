-- Add department enum type
CREATE TYPE "member_department" AS ENUM (
  'pengurus_harian',
  'kabid',
  'wasekbid',
  'wakil_bendahara'
);

-- Add department column to members table
ALTER TABLE "members" ADD COLUMN "department" member_department NOT NULL DEFAULT 'kabid';

-- Add index for department filtering
CREATE INDEX idx_members_department ON members(department);
CREATE INDEX idx_members_department_active ON members(department, is_active);
