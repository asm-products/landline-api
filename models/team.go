package models

import "time"

type Team struct {
	Id                string    `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Email             string    `json:"email"`
	EncryptedPassword string    `json:"encrypted_password"`
	Slug              string    `json:"slug"`
}
