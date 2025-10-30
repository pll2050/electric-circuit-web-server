-- Create users table for storing user information from Firebase OAuth
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(128) PRIMARY KEY,           -- Firebase UID (unique identifier)
    email VARCHAR(255) NOT NULL UNIQUE,    -- User email address
    display_name VARCHAR(255),             -- User display name
    photo_url TEXT,                        -- User profile photo URL
    provider VARCHAR(50) DEFAULT 'email',  -- Authentication provider (email, google, etc.)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);

-- Add comment to table
COMMENT ON TABLE users IS 'Stores user information from Firebase OAuth authentication';
COMMENT ON COLUMN users.id IS 'Firebase UID used as primary key';
COMMENT ON COLUMN users.email IS 'User email address (unique)';
COMMENT ON COLUMN users.display_name IS 'User display name from OAuth provider';
COMMENT ON COLUMN users.photo_url IS 'User profile photo URL from OAuth provider';
COMMENT ON COLUMN users.provider IS 'Authentication provider (email, google, github, etc.)';
COMMENT ON COLUMN users.created_at IS 'Timestamp when user was first registered';
COMMENT ON COLUMN users.updated_at IS 'Timestamp when user data was last updated';
