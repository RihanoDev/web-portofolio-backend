-- +goose Up
-- Junction table for posts and tags (many-to-many)
CREATE TABLE IF NOT EXISTS post_tags (
    post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, tag_id)
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX IF NOT EXISTS idx_post_tags_tag_id ON post_tags(tag_id);

-- +goose Down
DROP INDEX IF EXISTS idx_post_tags_tag_id;
DROP INDEX IF EXISTS idx_post_tags_post_id;
DROP TABLE IF EXISTS post_tags;
