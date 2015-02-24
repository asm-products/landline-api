-- +goose Up
CREATE TABLE nonces (
  id          		    uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  expires_at  		    timestamp   NOT NULL,
  nonce		            text        NOT NULL
);

CREATE UNIQUE INDEX idx_nonce on rooms (lower(slug));

-- +goose Down
DROP TABLE nonces;
