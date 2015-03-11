-- +goose Up
CREATE TABLE room_memberships (
  id          		    uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at  		    timestamp   NOT NULL,
  updated_at		      timestamp   NOT NULL,
  room_id             uuid        NOT NULL,
  user_id             uuid        NOT NULL,

  CONSTRAINT fk_room_memberships_rooms FOREIGN KEY (room_id) REFERENCES rooms (id),
  CONSTRAINT fk_room_memberships_users FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE UNIQUE INDEX idx_room_memberships_room_id_user_id on room_memberships (user_id, room_id);

-- +goose Down
DROP TABLE room_memberships;
