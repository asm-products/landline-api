package models

type SSORequest struct {
	Payload    string
	Signature  string
	Nonce      string
	TeamSlug   string
	ExternalId string
	AvatarUrl  string
	Username   string
	Email      string
	ProfileUrl string
	RealName   string
}

func (r *SSORequest) IsValid(secret string) bool {
	return Sign([]byte(secret), []byte(r.Payload)) == r.Signature
}
