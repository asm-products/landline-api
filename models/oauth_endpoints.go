package models

type OAuthEndpoints struct {
	AuthorizeUrl string `db:"oauth_authorize_url" json:"authorize_url"`
	TokenUrl     string `db:"oauth_token_url"     json:"token_url"`
}
