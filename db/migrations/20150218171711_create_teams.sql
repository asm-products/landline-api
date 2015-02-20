-- +goose Up
CREATE TABLE teams (
  id          		    uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at  		    timestamp   NOT NULL,
  updated_at		      timestamp   NOT NULL,
  email       		    text        NOT NULL,
  encrypted_password	text 	      NOT NULL,
  oauth_authorize_url	text 	      NOT NULL,
  oauth_token_url	    text 	      NOT NULL,
  slug                text 	      NOT NULL
);

CREATE UNIQUE INDEX idx_teams_email on teams (lower(email));

-- +goose Down
DROP TABLE teams;
