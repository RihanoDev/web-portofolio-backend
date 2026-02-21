-- +goose Up
-- Junction and related tables for projects

-- Project-Tag junction table (many-to-many)
CREATE TABLE IF NOT EXISTS project_tags (
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, tag_id)
);

-- Project-Technology junction table (technologies are tags)
CREATE TABLE IF NOT EXISTS project_technologies (
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, tag_id)
);

-- Project images
CREATE TABLE IF NOT EXISTS project_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    caption TEXT,
    alt_text TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Project videos
CREATE TABLE IF NOT EXISTS project_videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    caption TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_project_tags_project_id ON project_tags(project_id);
CREATE INDEX IF NOT EXISTS idx_project_tags_tag_id ON project_tags(tag_id);
CREATE INDEX IF NOT EXISTS idx_project_technologies_project_id ON project_technologies(project_id);
CREATE INDEX IF NOT EXISTS idx_project_technologies_tag_id ON project_technologies(tag_id);
CREATE INDEX IF NOT EXISTS idx_project_images_project_id ON project_images(project_id);
CREATE INDEX IF NOT EXISTS idx_project_videos_project_id ON project_videos(project_id);

-- +goose Down
DROP INDEX IF EXISTS idx_project_videos_project_id;
DROP INDEX IF EXISTS idx_project_images_project_id;
DROP INDEX IF EXISTS idx_project_technologies_tag_id;
DROP INDEX IF EXISTS idx_project_technologies_project_id;
DROP INDEX IF EXISTS idx_project_tags_tag_id;
DROP INDEX IF EXISTS idx_project_tags_project_id;
DROP TABLE IF EXISTS project_videos;
DROP TABLE IF EXISTS project_images;
DROP TABLE IF EXISTS project_technologies;
DROP TABLE IF EXISTS project_tags;
