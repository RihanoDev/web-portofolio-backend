-- +goose Up
-- Seed sample categories (if they don't exist)
INSERT INTO categories (name, slug, description) VALUES 
('Web Development', 'web-development', 'Projects related to web development')
ON CONFLICT (slug) DO NOTHING;

INSERT INTO categories (name, slug, description) VALUES 
('Mobile Development', 'mobile-development', 'Projects related to mobile app development')
ON CONFLICT (slug) DO NOTHING;

INSERT INTO categories (name, slug, description) VALUES 
('DevOps', 'devops', 'Projects related to DevOps and infrastructure')
ON CONFLICT (slug) DO NOTHING;

-- Seed sample tags (if they don't exist)
INSERT INTO tags (name, slug) VALUES 
('React', 'react'),
('TypeScript', 'typescript'),
('Go', 'go'),
('PostgreSQL', 'postgresql'),
('Docker', 'docker'),
('Vue.js', 'vuejs'),
('Node.js', 'nodejs'),
('Python', 'python')
ON CONFLICT (slug) DO NOTHING;

-- Seed sample projects (using admin user ID from previous migration)
INSERT INTO projects (id, title, slug, description, content, thumbnail_url, status, category_id, author_id, github_url, live_demo_url) 
SELECT 
    gen_random_uuid(),
    'Portfolio CMS System',
    'portfolio-cms-system',
    'A comprehensive content management system for personal portfolios built with Go backend and React frontend.',
    '<h2>Portfolio CMS System</h2><p>This is a full-stack portfolio management system that allows users to manage their personal portfolio content including projects, articles, and experiences.</p><h3>Features</h3><ul><li>User authentication and authorization</li><li>Project management with image and video support</li><li>Article publishing system</li><li>Experience timeline management</li><li>Analytics and view tracking</li><li>Responsive design</li></ul><h3>Technology Stack</h3><ul><li>Backend: Go with Gin framework</li><li>Frontend: React with TypeScript</li><li>Database: PostgreSQL</li><li>Authentication: JWT</li><li>Containerization: Docker</li></ul>',
    'https://via.placeholder.com/400x300?text=Portfolio+CMS',
    'published',
    c.id,
    u.id,
    'https://github.com/username/portfolio-cms',
    'https://portfolio-cms.example.com'
FROM categories c, users u 
WHERE c.slug = 'web-development' AND u.email = 'admin@example.com'
AND NOT EXISTS (SELECT 1 FROM projects WHERE slug = 'portfolio-cms-system');

INSERT INTO projects (id, title, slug, description, content, thumbnail_url, status, category_id, author_id, github_url, live_demo_url) 
SELECT 
    gen_random_uuid(),
    'Task Management API',
    'task-management-api',
    'A RESTful API for task management with real-time notifications.',
    '<h2>Task Management API</h2><p>A robust REST API built with Go for managing tasks and projects with real-time updates.</p><h3>Key Features</h3><ul><li>RESTful API design</li><li>Real-time notifications</li><li>User collaboration</li><li>Task categorization</li><li>Progress tracking</li></ul>',
    'https://via.placeholder.com/400x300?text=Task+API',
    'published',
    c.id,
    u.id,
    'https://github.com/username/task-api',
    NULL
FROM categories c, users u 
WHERE c.slug = 'web-development' AND u.email = 'admin@example.com'
AND NOT EXISTS (SELECT 1 FROM projects WHERE slug = 'task-management-api');

-- Seed sample articles
INSERT INTO articles (id, title, slug, excerpt, content, featured_image_url, status, author_id) 
SELECT 
    gen_random_uuid(),
    'Building Scalable Web Applications with Go',
    'building-scalable-web-applications-go',
    'Learn how to build high-performance, scalable web applications using Go programming language.',
    '<h1>Building Scalable Web Applications with Go</h1><p>Go (Golang) has become increasingly popular for building web applications due to its simplicity, performance, and excellent concurrency support.</p><h2>Why Choose Go for Web Development?</h2><ul><li>Excellent performance characteristics</li><li>Built-in concurrency with goroutines</li><li>Strong standard library</li><li>Fast compilation times</li><li>Great for microservices</li></ul><h2>Getting Started</h2><p>In this article, we''ll explore the fundamentals of building web applications with Go...</p>',
    'https://via.placeholder.com/800x400?text=Go+Web+Development',
    'published',
    u.id
FROM users u 
WHERE u.email = 'admin@example.com'
AND NOT EXISTS (SELECT 1 FROM articles WHERE slug = 'building-scalable-web-applications-go');

INSERT INTO articles (id, title, slug, excerpt, content, featured_image_url, status, author_id) 
SELECT 
    gen_random_uuid(),
    'Modern Frontend Development with React and TypeScript',
    'modern-frontend-development-react-typescript',
    'Explore best practices for building modern frontend applications with React and TypeScript.',
    '<h1>Modern Frontend Development with React and TypeScript</h1><p>React combined with TypeScript provides a powerful foundation for building robust frontend applications.</p><h2>Benefits of TypeScript</h2><ul><li>Static type checking</li><li>Better IDE support</li><li>Improved code maintainability</li><li>Enhanced developer experience</li></ul><p>Let''s dive into the best practices...</p>',
    'https://via.placeholder.com/800x400?text=React+TypeScript',
    'published',
    u.id
FROM users u 
WHERE u.email = 'admin@example.com'
AND NOT EXISTS (SELECT 1 FROM articles WHERE slug = 'modern-frontend-development-react-typescript');

-- Seed sample experiences
INSERT INTO experiences (title, company, location, start_date, end_date, current, description, responsibilities, technologies, company_url) 
SELECT 
    'Senior Full Stack Developer',
    'Tech Solutions Inc.',
    'Jakarta, Indonesia',
    '2022-01-01',
    NULL,
    true,
    'Leading the development of scalable web applications and mentoring junior developers.',
    to_jsonb(ARRAY['Lead development team of 5 developers', 'Architect and implement new features', 'Code review and quality assurance', 'Collaborate with product and design teams', 'Optimize application performance']),
    to_jsonb(ARRAY['React', 'TypeScript', 'Node.js', 'PostgreSQL', 'Docker', 'AWS']),
    'https://techsolutions.example.com'
WHERE NOT EXISTS (SELECT 1 FROM experiences WHERE company = 'Tech Solutions Inc.' AND title = 'Senior Full Stack Developer');

INSERT INTO experiences (title, company, location, start_date, end_date, current, description, responsibilities, technologies, company_url) 
SELECT 
    'Full Stack Developer',
    'Digital Agency Pro',
    'Bandung, Indonesia',
    '2020-03-15',
    '2021-12-31',
    false,
    'Developed and maintained multiple client projects using modern web technologies.',
    to_jsonb(ARRAY['Develop responsive web applications', 'Integrate third-party APIs', 'Database design and optimization', 'Deploy applications to cloud platforms', 'Client communication and requirement gathering']),
    to_jsonb(ARRAY['Vue.js', 'PHP', 'Laravel', 'MySQL', 'JavaScript', 'CSS']),
    'https://digitalagencypro.example.com'
WHERE NOT EXISTS (SELECT 1 FROM experiences WHERE company = 'Digital Agency Pro' AND title = 'Full Stack Developer');

INSERT INTO experiences (title, company, location, start_date, end_date, current, description, responsibilities, technologies, company_url) 
SELECT 
    'Junior Software Developer',
    'StartupXYZ',
    'Yogyakarta, Indonesia',
    '2019-06-01',
    '2020-02-28',
    false,
    'Started my career as a junior developer working on various web development projects.',
    to_jsonb(ARRAY['Implement frontend components', 'Write unit tests', 'Bug fixing and maintenance', 'Participate in agile development process', 'Learn new technologies and frameworks']),
    to_jsonb(ARRAY['HTML', 'CSS', 'JavaScript', 'Python', 'Django', 'Git']),
    'https://startupxyz.example.com'
WHERE NOT EXISTS (SELECT 1 FROM experiences WHERE company = 'StartupXYZ' AND title = 'Junior Software Developer');

-- Create associations for projects and tags
INSERT INTO project_technologies (project_id, tag_id) 
SELECT p.id, t.id 
FROM projects p, tags t 
WHERE p.slug = 'portfolio-cms-system' AND t.slug IN ('react', 'typescript', 'go', 'postgresql', 'docker');

INSERT INTO project_technologies (project_id, tag_id) 
SELECT p.id, t.id 
FROM projects p, tags t 
WHERE p.slug = 'task-management-api' AND t.slug IN ('go', 'postgresql', 'docker');

-- Create associations for articles and tags  
INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id 
FROM articles a, tags t 
WHERE a.slug = 'building-scalable-web-applications-go' AND t.slug IN ('go', 'nodejs');

INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id 
FROM articles a, tags t 
WHERE a.slug = 'modern-frontend-development-react-typescript' AND t.slug IN ('react', 'typescript');

-- +goose Down
-- Remove seeded data
DELETE FROM project_technologies WHERE project_id IN (
    SELECT id FROM projects WHERE slug IN ('portfolio-cms-system', 'task-management-api')
);

DELETE FROM article_tags WHERE article_id IN (
    SELECT id FROM articles WHERE slug IN ('building-scalable-web-applications-go', 'modern-frontend-development-react-typescript')
);

DELETE FROM experiences WHERE company IN ('Tech Solutions Inc.', 'Digital Agency Pro', 'StartupXYZ');
DELETE FROM articles WHERE slug IN ('building-scalable-web-applications-go', 'modern-frontend-development-react-typescript');
DELETE FROM projects WHERE slug IN ('portfolio-cms-system', 'task-management-api');
DELETE FROM tags WHERE slug IN ('react', 'typescript', 'go', 'postgresql', 'docker', 'vuejs', 'nodejs', 'python');
DELETE FROM categories WHERE slug IN ('web-development', 'mobile-development', 'devops');
