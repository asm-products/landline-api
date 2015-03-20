
-- +goose Up
DROP INDEX idx_nonce;


-- +goose Down
CREATE UNIQUE INDEX idx_nonce on rooms (lower(slug));
