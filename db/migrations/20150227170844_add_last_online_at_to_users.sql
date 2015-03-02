
-- +goose Up
ALTER TABLE users ADD COLUMN last_online_at timestamp;
UPDATE users SET last_online_at = created_at;

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS last_online;
