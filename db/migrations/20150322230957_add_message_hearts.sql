
-- +goose Up
CREATE TABLE message_hearts (
	message_id          uuid        NOT NULL,
	user_id             uuid        NOT NULL,
	created_at  		    timestamp   NOT NULL,
	PRIMARY KEY (message_id, user_id)
);

-- +goose Down
DROP TABLE message_hearts;
