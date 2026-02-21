-- +goose Up
-- Create combined index on page_views for common analytics queries
CREATE INDEX IF NOT EXISTS idx_page_views_page_timestamp ON page_views (page, timestamp);
CREATE INDEX IF NOT EXISTS idx_page_views_visitor_timestamp ON page_views (visitor_id, timestamp);

-- +goose Down
DROP INDEX IF EXISTS idx_page_views_page_timestamp;
DROP INDEX IF EXISTS idx_page_views_visitor_timestamp;
