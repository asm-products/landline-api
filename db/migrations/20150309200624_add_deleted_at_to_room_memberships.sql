
-- +goose Up
ALTER TABLE room_memberships ADD COLUMN deleted_at timestamp;

-- +goose Down
ALTER TABLE room_memberships DROP COLUMN IF EXISTS deleted_at;
