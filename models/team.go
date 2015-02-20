package models

import "time"

type Team struct {
	Id                string    `db:"id" 									json:"id"`
	CreatedAt         time.Time `db:"created_at" 					json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" 					json:"updated_at"`
	Email             string    `db:"email" 							json:"email"`
	EncryptedPassword string    `db:"encrypted_password" 	json:"encrypted_password"`
	OAuthAuthorizeUrl string    `db:"oauth_authorize_url" json:"oauth_authorize_url"`
	OAuthTokenUrl     string    `db:"oauth_token_url" 		json:"oauth_token_url"`
	Slug              string    `db:"slug" 								json:"slug"`
}
