
-- +goose Up
CREATE UNIQUE INDEX idx_rooms_slug_team_id on rooms (lower(slug), team_id);


-- +goose Down
DROP INDEX idx_rooms_slug_team_id;
