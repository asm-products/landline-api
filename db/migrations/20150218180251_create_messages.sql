-- +goose Up
CREATE TABLE messages (
  id          		    uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at  		    timestamp   NOT NULL,
  updated_at		      timestamp   NOT NULL,
  room_id             uuid        NOT NULL,
  user_id             uuid        NOT NULL,
  body                text        NOT NULL,

  CONSTRAINT fk_messages_rooms FOREIGN KEY (room_id) REFERENCES rooms (id),
  CONSTRAINT fk_messages_users FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE messages;
