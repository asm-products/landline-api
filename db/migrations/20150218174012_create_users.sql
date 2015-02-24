-- +goose Up
CREATE TABLE users (
  id          		    uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at  		    timestamp   NOT NULL,
  updated_at		      timestamp   NOT NULL,
  team_id             uuid        NOT NULL,
  avatar_url	        text,
  email         		  text        NOT NULL,
  external_id         text        NOT NULL,
  profile_url         text,
  real_name           text,
  username       		  text        NOT NULL,

  CONSTRAINT fk_users_teams FOREIGN KEY (team_id) REFERENCES teams (id)
);

CREATE UNIQUE INDEX idx_users_username on users (lower(username));
CREATE UNIQUE INDEX idx_users_external_id on users (lower(external_id));

-- +goose Down
DROP TABLE users;
