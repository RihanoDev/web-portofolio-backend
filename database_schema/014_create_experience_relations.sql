-- +goose Up
-- Junction table for experiences

-- Experience-Technology junction table (technologies are tags)
CREATE TABLE IF NOT EXISTS experience_technologies (
    experience_id INTEGER REFERENCES experiences(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (experience_id, tag_id)
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_experience_technologies_experience_id ON experience_technologies(experience_id);
CREATE INDEX IF NOT EXISTS idx_experience_technologies_tag_id ON experience_technologies(tag_id);

-- +goose Down
DROP INDEX IF EXISTS idx_experience_technologies_tag_id;
DROP INDEX IF EXISTS idx_experience_technologies_experience_id;
DROP TABLE IF EXISTS experience_technologies;
