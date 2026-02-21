-- Create junction tables for many-to-many relationships

-- Project-Tag junction table (GORM expects this name)
CREATE TABLE IF NOT EXISTS project_tags (
    project_id UUID NOT NULL,
    tag_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, tag_id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Experience-Technology junction table (technologies are basically tags)
CREATE TABLE IF NOT EXISTS experience_technologies (
    experience_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (experience_id, tag_id),
    FOREIGN KEY (experience_id) REFERENCES experiences(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Add indexes for better performance
CREATE INDEX IF NOT EXISTS idx_project_tags_project_id ON project_tags(project_id);
CREATE INDEX IF NOT EXISTS idx_project_tags_tag_id ON project_tags(tag_id);
CREATE INDEX IF NOT EXISTS idx_experience_technologies_experience_id ON experience_technologies(experience_id);
CREATE INDEX IF NOT EXISTS idx_experience_technologies_tag_id ON experience_technologies(tag_id);
