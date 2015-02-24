package models

import (
	"database/sql"
	"time"

	"gopkg.in/gorp.v1"
)

type User struct {
	Id        string    `db:"id" 					json:"id"`
	CreatedAt time.Time `db:"created_at" 	json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" 	json:"updated_at"`
	TeamId    string    `db:"team_id" 		json:"team_id"`

	AvatarUrl  string `db:"avatar_url" 		json:"avatar_url"`
	Email      string `db:"email" 				json:"email"`
	ExternalId string `db:"external_id" 	json:"external_id"`
	ProfileUrl string `db:"profile_url" 	json:"profile_url"`
	RealName   string `db:"real_name" 		json:"real_name"`
	Username   string `db:"username" 			json:"username"`
}

func FindOrCreateUserByExternalId(fields *User) (*User, error) {
	var user User
	err := Db.SelectOne(&user, `select * from users where external_id = $1 limit 1`, fields.ExternalId)
	if err == sql.ErrNoRows {
		err = Db.Insert(fields)
		return fields, err
	}
	return &user, err
}

func FindUser(id string) (*User, error) {
	var user User
	err := Db.SelectOne(&user, `select * from users where id = $1 limit 1`, id)
	return &user, err
}

func (o *User) PreInsert(s gorp.SqlExecutor) error {
	o.CreatedAt = time.Now()
	o.UpdatedAt = o.CreatedAt
	return nil
}

func (o *User) PreUpdate(s gorp.SqlExecutor) error {
	o.UpdatedAt = time.Now()
	return nil
}
