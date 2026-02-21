-- +goose Up
-- Add excerpt column to posts table
ALTER TABLE posts ADD COLUMN excerpt TEXT;

-- +goose Down
-- Remove the excerpt column
ALTER TABLE posts DROP COLUMN excerpt;
