
-- +goose Up
DROP INDEX idx_users_external_id;


-- +goose Down
CREATE UNIQUE INDEX idx_users_external_id on users (lower(external_id));
