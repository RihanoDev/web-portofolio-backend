-- +goose Up
CREATE TABLE IF NOT EXISTS content_views (
    id SERIAL PRIMARY KEY,
    content_id VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    visitor_id VARCHAR(255) NOT NULL,
    user_agent TEXT,
    referrer TEXT,
    ip VARCHAR(100),
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_content_views_content ON content_views(content_id, type);
CREATE INDEX IF NOT EXISTS idx_content_views_timestamp ON content_views(timestamp);
CREATE INDEX IF NOT EXISTS idx_content_views_visitor ON content_views(visitor_id);

-- +goose Down
DROP TABLE IF EXISTS content_views;
