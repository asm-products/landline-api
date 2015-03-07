
-- +goose Up
ALTER TABLE teams ADD COLUMN webhook_url text;

-- +goose Down
ALTER TABLE teams DROP COLUMN IF EXISTS webhook_url;
