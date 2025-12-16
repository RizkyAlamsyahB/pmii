-- Remove indexes
DROP INDEX IF EXISTS idx_members_department_active;
DROP INDEX IF EXISTS idx_members_department;

-- Remove department column
ALTER TABLE "members" DROP COLUMN IF EXISTS "department";

-- Remove department enum type
DROP TYPE IF EXISTS "member_department";
