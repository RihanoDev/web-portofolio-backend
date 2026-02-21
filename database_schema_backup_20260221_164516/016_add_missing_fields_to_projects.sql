-- +goose Up
-- +goose StatementBegin
ALTER TABLE projects
ADD COLUMN IF NOT EXISTS git_hub_url VARCHAR(255),
ADD COLUMN IF NOT EXISTS deployment_url VARCHAR(255),
ADD COLUMN IF NOT EXISTS start_date TIMESTAMP,
ADD COLUMN IF NOT EXISTS end_date TIMESTAMP,
ADD COLUMN IF NOT EXISTS is_featured BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS order_number INT DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE projects
DROP COLUMN IF EXISTS git_hub_url,
DROP COLUMN IF EXISTS deployment_url,
DROP COLUMN IF EXISTS start_date,
DROP COLUMN IF EXISTS end_date,
DROP COLUMN IF EXISTS is_featured,
DROP COLUMN IF EXISTS order_number;
-- +goose StatementEnd
