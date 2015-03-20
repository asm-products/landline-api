
-- +goose Up
CREATE UNIQUE INDEX idx_nonce_room_team_id on rooms (lower(slug), team_id);


-- +goose Down
DROP INDEX idx_nonce_room_team_id;
