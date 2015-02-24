package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gopkg.in/gorp.v1"
)

type Team struct {
	Id                string    `db:"id" 									json:"id"`
	CreatedAt         time.Time `db:"created_at" 					json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" 					json:"updated_at"`
	Email             string    `db:"email" 							json:"email"`
	EncryptedPassword string    `db:"encrypted_password" 	json:"encrypted_password"`
	SSOSecret         string    `db:"sso_secret" 					json:"-"`
	SSOUrl            string    `db:"sso_url" 						json:"sso_url"`
	Slug              string    `db:"slug" 								json:"slug"`
}

func FindTeamBySlug(slug string) *Team {
	var team Team
	err := Db.SelectOne(&team, "select * from Teams where slug=$1", slug)
	if err != nil {
		panic(err)
	}
	return &team
}

func ShaString(raw []byte) string {
	hasher := sha256.New()
	hasher.Write(raw)
	return hex.EncodeToString(hasher.Sum(nil))
}

func Sign(secret, payload []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func (o *Team) PreInsert(s gorp.SqlExecutor) error {
	o.CreatedAt = time.Now()
	o.UpdatedAt = o.CreatedAt
	return nil
}

func (o *Team) PreUpdate(s gorp.SqlExecutor) error {
	o.UpdatedAt = time.Now()
	return nil
}
