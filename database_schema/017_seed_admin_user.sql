-- +goose Up
-- Insert admin user with hashed password
INSERT INTO users (username, email, password_hash, role, bio, is_active) 
VALUES (
    'admin', 
    'admin@portfolio.com', 
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: password
    'admin',
    'Portfolio Administrator',
    true
) ON CONFLICT (username) DO NOTHING;

-- Insert default editor user
INSERT INTO users (username, email, password_hash, role, bio, is_active) 
VALUES (
    'editor', 
    'editor@portfolio.com', 
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: password
    'editor',
    'Content Editor',
    true
) ON CONFLICT (username) DO NOTHING;

-- +goose Down
DELETE FROM users WHERE username IN ('admin', 'editor');
