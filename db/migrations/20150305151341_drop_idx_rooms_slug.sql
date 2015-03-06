
-- +goose Up
DROP INDEX idx_rooms_slug;


-- +goose Down
CREATE UNIQUE INDEX idx_rooms_slug on rooms (lower(slug));
