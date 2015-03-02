
-- +goose Up
CREATE UNIQUE INDEX idx_users_external_id_team_id on users (lower(external_id), team_id);


-- +goose Down
DROP INDEX idx_users_external_id_team_id;
