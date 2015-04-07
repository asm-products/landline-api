package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

const (
	NonceBytes = 32
)

type Nonce struct {
	Id        string    `db:"id" json:"-"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	Nonce     string    `db:"nonce" json:"nonce"`
}

func CreateNonce() (*Nonce, error) {
	n, err := generate()
	if err != nil {
		panic(err)
	}
	nonce := &Nonce{
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Nonce:     n,
	}
	err = Db.Insert(nonce)
	return nonce, err
}

func generate() (string, error) {
	b := make([]byte, NonceBytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	str := hex.EncodeToString(b)
	return string(str[:NonceBytes]), nil
}

func NonceValid(nonce string) bool {
	id, err := Db.SelectNullStr(
		"select id from nonces where nonce = $1 and expires_at > $2 limit 1",
		nonce, time.Now(),
	)
	if err != nil {
		panic(err)
	}
	if !id.Valid {
		return false
	}

	_, err = Db.Exec("delete from nonces where id=$1", id.String)
	if err != nil {
		panic(err)
	}

	return true
}
