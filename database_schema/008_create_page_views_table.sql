-- +goose Up
CREATE TABLE IF NOT EXISTS page_views (
    id SERIAL PRIMARY KEY,
    page VARCHAR(255) NOT NULL,
    visitor_id VARCHAR(255) NOT NULL,
    user_agent TEXT,
    referrer TEXT,
    ip VARCHAR(100),
    country VARCHAR(100),
    city VARCHAR(100),
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_page_views_page ON page_views(page);
CREATE INDEX IF NOT EXISTS idx_page_views_timestamp ON page_views(timestamp);
CREATE INDEX IF NOT EXISTS idx_page_views_visitor ON page_views(visitor_id);

-- +goose Down
DROP TABLE IF EXISTS page_views;
