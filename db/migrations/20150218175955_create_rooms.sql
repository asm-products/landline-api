-- +goose Up
CREATE TABLE rooms (
  id          		    uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at  		    timestamp   NOT NULL,
  updated_at		      timestamp   NOT NULL,
  team_id             uuid        NOT NULL,
  slug                text        NOT NULL,
  topic               text        NOT NULL,

  CONSTRAINT fk_rooms_teams FOREIGN KEY (team_id) REFERENCES teams (id)
);

CREATE UNIQUE INDEX idx_rooms_slug on rooms (lower(slug));

-- +goose Down
DROP TABLE rooms;
