package models

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/gorp.v1"
)

type User struct {
	Id           string    `db:"id" json:"id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	LastOnlineAt time.Time `db:"last_online_at" json:"last_online_at"`
	TeamId       string    `db:"team_id" json:"team_id"`

	AvatarUrl  string `db:"avatar_url" json:"avatar_url"`
	Email      string `db:"email" json:"email"`
	ExternalId string `db:"external_id" json:"external_id"`
	ProfileUrl string `db:"profile_url" json:"profile_url"`
	RealName   string `db:"real_name" json:"real_name"`
	Username   string `db:"username" json:"username"`
}

func FindOrCreateUserByExternalId(fields *User) (*User, error) {
	var user User
	err := Db.SelectOne(&user, `select * from users where external_id = $1 limit 1`, fields.ExternalId)
	if err == sql.ErrNoRows {
		err = Db.Insert(fields)
		return fields, err
	}
	_, err = Db.Update(&user)

	return &user, err
}

func FindUser(id string) (*User, error) {
	var user User
	err := Db.SelectOne(&user, `select * from users where id = $1 limit 1`, id)
	return &user, err
}

func FindUsers(teamId string) ([]User, error) {
	var users []User
	_, err := Db.Select(
		&users,
		`SELECT * FROM users WHERE team_id = $1`,
		teamId,
	)

	return users, err
}

func FindRecentlyOnlineUsers(teamId string) ([]User, error) {
	var users []User
	_, err := Db.Select(
		&users,
		`SELECT * FROM users WHERE team_id = $1
		and last_online_at >= now() - '2 hour'::INTERVAL`,
		teamId,
	)

	return users, err
}

func UnreadRooms(userId string) ([]byte, error) {
	req, err := http.NewRequest(
		"GET",
		os.Getenv("RR_URL")+"/readers/"+userId,
		nil,
	)

	req.SetBasicAuth(os.Getenv("RR_PRIVATE_KEY"), "")

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func (o *User) PreInsert(s gorp.SqlExecutor) error {
	o.CreatedAt = time.Now()
	o.UpdatedAt = o.CreatedAt
	o.LastOnlineAt = o.CreatedAt
	return nil
}

func (o *User) PreUpdate(s gorp.SqlExecutor) error {
	o.UpdatedAt = time.Now()
	o.LastOnlineAt = o.UpdatedAt
	return nil
}
