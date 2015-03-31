package models

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/gorp.v1"
)

type Room struct {
	Id        string     `db:"id" json:"id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	TeamId    string     `db:"team_id" json:"team_id"`
	Slug      string     `db:"slug" json:"slug"`
	Topic     string     `db:"topic" json:"topic"`
}

func DeleteRoom(slug, teamId string) error {
	var room Room
	err := Db.SelectOne(&room, "select * from rooms where slug=$1 and team_id=$2", slug, teamId)
	if err != nil {
		panic(err)
	}
	t := time.Now()
	room.DeletedAt = &t
	_, err = Db.Update(room)
	return err
}

func FindOrCreateRoom(fields *Room) (*Room, error) {
	var room Room
	err := Db.SelectOne(&room, "select * from rooms where slug=$1 and team_id=$2", fields.Slug, fields.TeamId)
	if err == sql.ErrNoRows {
		err = Db.Insert(fields)
	}
	return &room, err
}

func FindRoom(slug, teamId string) (*Room, error) {
	var room Room
	err := Db.SelectOne(&room, "select * from rooms where slug=$1 and team_id=$2", slug, teamId)
	return &room, err
}

func FindRooms(teamId string) ([]Room, error) {
	var rooms []Room
	_, err := Db.Select(
		&rooms,
		`select * from rooms where team_id = $1`,
		teamId,
	)

	return rooms, err
}

func FindRoomById(id string) *Room {
	var room Room
	err := Db.SelectOne(&room, "select * from rooms where id = $1", id)

	if err != nil {
		panic(err)
	}

	return &room
}

func Subscribers(roomId string) (*[]string, error) {
	var subscribers []string

	_, err := Db.Select(
		&subscribers,
		`select user_id from room_memberships where room_id = $1`,
		roomId,
	)

	if err == sql.ErrNoRows {
		return &subscribers, nil
	}

	return &subscribers, err
}

func UnreadRooms(userId string) ([]byte, error) {
	req, err := http.NewRequest(
		"GET",
		os.Getenv("RR_URL")+"/readers/"+userId,
		nil,
	)

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(os.Getenv("RR_PRIVATE_KEY"), "")

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func UpdateRoom(slug, teamId string, fields *Room) (*Room, error) {
	var room Room
	err := Db.SelectOne(
		&room,
		"select * from rooms where slug=$1 and team_id=$2",
		slug,
		teamId,
	)
	if err != nil {
		panic(err)
	}

	room.Slug = fields.Slug
	room.Topic = fields.Topic

	_, err = Db.Update(&room)

	return &room, err
}

func (o *Room) PreInsert(s gorp.SqlExecutor) error {
	o.CreatedAt = time.Now()
	o.UpdatedAt = o.CreatedAt
	return nil
}

func (o *Room) PreUpdate(s gorp.SqlExecutor) error {
	o.UpdatedAt = time.Now()
	return nil
}
