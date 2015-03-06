
-- +goose Up
ALTER TABLE rooms ADD COLUMN deleted_at timestamp;


-- +goose Down
ALTER TABLE rooms DROP COLUMN deleted_at;
