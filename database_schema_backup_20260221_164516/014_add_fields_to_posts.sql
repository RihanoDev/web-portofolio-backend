-- +goose Up
-- Create the UUID extension if it doesn't exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Add missing fields to posts table to match the model
-- excerpt column already added in migration 013
ALTER TABLE posts ADD COLUMN IF NOT EXISTS featured_image_url TEXT;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS read_time INTEGER DEFAULT 0;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS view_count INTEGER DEFAULT 0;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS metadata JSONB;

-- +goose Down
-- Revert changes if needed
-- excerpt handled in migration 013
ALTER TABLE posts DROP COLUMN IF EXISTS featured_image_url;
ALTER TABLE posts DROP COLUMN IF EXISTS read_time;
ALTER TABLE posts DROP COLUMN IF EXISTS view_count;
ALTER TABLE posts DROP COLUMN IF EXISTS metadata;
