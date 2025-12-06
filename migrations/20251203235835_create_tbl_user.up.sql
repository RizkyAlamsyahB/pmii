-- Create tbl_user table
-- This table stores user information for the PMII CMS
-- Matches legacy database structure for backward compatibility

CREATE TABLE IF NOT EXISTS tbl_user (
    user_id SERIAL PRIMARY KEY,
    user_name VARCHAR(100) NOT NULL,
    user_email VARCHAR(60) NOT NULL,
    user_password VARCHAR(255) NOT NULL, -- Changed from VARCHAR(40) to support bcrypt hashes
    user_level VARCHAR(10) NOT NULL DEFAULT '2', -- 1=Admin, 2=User
    user_status VARCHAR(10) NOT NULL DEFAULT '1', -- 1=Active, 0=Inactive
    user_photo VARCHAR(40) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_user_email ON tbl_user(user_email);
CREATE INDEX idx_user_level ON tbl_user(user_level);
CREATE INDEX idx_user_status ON tbl_user(user_status);

-- Add unique constraint on email
ALTER TABLE tbl_user ADD CONSTRAINT unique_user_email UNIQUE (user_email);

-- Add comments for documentation
COMMENT ON TABLE tbl_user IS 'User table for PMII CMS authentication and authorization';
COMMENT ON COLUMN tbl_user.user_id IS 'Primary key, auto-incrementing user ID';
COMMENT ON COLUMN tbl_user.user_name IS 'Full name of the user';
COMMENT ON COLUMN tbl_user.user_email IS 'Unique email address for login';
COMMENT ON COLUMN tbl_user.user_password IS 'Bcrypt hashed password (migrated from legacy MD5)';
COMMENT ON COLUMN tbl_user.user_level IS 'User role: 1=Admin, 2=Author, 3=Contributor';
COMMENT ON COLUMN tbl_user.user_status IS 'Account status: 1=Active, 0=Inactive';
COMMENT ON COLUMN tbl_user.user_photo IS 'Filename of user profile photo (stored in storage)';
