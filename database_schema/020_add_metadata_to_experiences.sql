-- +goose Up
-- Add metadata column to experiences table
ALTER TABLE experiences ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}'::jsonb;

-- +goose Down
-- Remove metadata column from experiences table
ALTER TABLE experiences DROP COLUMN IF EXISTS metadata;
