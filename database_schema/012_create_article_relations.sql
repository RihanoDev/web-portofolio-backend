-- +goose Up
-- Junction tables for articles

-- Article-Category junction table (many-to-many)
CREATE TABLE IF NOT EXISTS article_categories (
    article_id UUID REFERENCES articles(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (article_id, category_id)
);

-- Article-Tag junction table (many-to-many)
CREATE TABLE IF NOT EXISTS article_tags (
    article_id UUID REFERENCES articles(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (article_id, tag_id)
);

-- Article images
CREATE TABLE IF NOT EXISTS article_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id UUID REFERENCES articles(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    caption TEXT,
    alt_text TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Article videos
CREATE TABLE IF NOT EXISTS article_videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id UUID REFERENCES articles(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    caption TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_article_categories_article_id ON article_categories(article_id);
CREATE INDEX IF NOT EXISTS idx_article_categories_category_id ON article_categories(category_id);
CREATE INDEX IF NOT EXISTS idx_article_tags_article_id ON article_tags(article_id);
CREATE INDEX IF NOT EXISTS idx_article_tags_tag_id ON article_tags(tag_id);
CREATE INDEX IF NOT EXISTS idx_article_images_article_id ON article_images(article_id);
CREATE INDEX IF NOT EXISTS idx_article_videos_article_id ON article_videos(article_id);

-- +goose Down
DROP INDEX IF EXISTS idx_article_videos_article_id;
DROP INDEX IF EXISTS idx_article_images_article_id;
DROP INDEX IF EXISTS idx_article_tags_tag_id;
DROP INDEX IF NOT EXISTS idx_article_tags_article_id;
DROP INDEX IF EXISTS idx_article_categories_category_id;
DROP INDEX IF EXISTS idx_article_categories_article_id;
DROP TABLE IF EXISTS article_videos;
DROP TABLE IF EXISTS article_images;
DROP TABLE IF EXISTS article_tags;
DROP TABLE IF EXISTS article_categories;
