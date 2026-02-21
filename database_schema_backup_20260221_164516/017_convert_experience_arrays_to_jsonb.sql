-- +goose Up
-- Convert array columns to JSONB for better compatibility
ALTER TABLE experiences 
ALTER COLUMN responsibilities TYPE jsonb USING to_jsonb(responsibilities),
ALTER COLUMN technologies TYPE jsonb USING to_jsonb(technologies);

-- +goose Down  
-- Revert back to text array (this might lose data formatting)
ALTER TABLE experiences 
ALTER COLUMN responsibilities TYPE text[] USING ARRAY(SELECT jsonb_array_elements_text(responsibilities)),
ALTER COLUMN technologies TYPE text[] USING ARRAY(SELECT jsonb_array_elements_text(technologies));
