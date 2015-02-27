package models

import (
	"database/sql"
	"time"
)

type Room struct {
	Id        string    `db:"id" 				 json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	TeamId    string    `db:"team_id" 	 json:"team_id"`
	Slug      string    `db:"slug" 			 json:"slug"`
	Topic     string    `db:"topic" 		 json:"topic"`
}

func FindOrCreateRoom(fields *Room) (*Room, error) {
	var room Room
	err := Db.SelectOne(&room, "select * from Rooms where slug=$1 and team_id=$2", fields.Slug, fields.TeamId)
	if err == sql.ErrNoRows {
		err = Db.Insert(fields)
		return fields, err
	}
	return &room, err
}

func FindRoom(slug string, teamId string) (*Room, error) {
	var room Room
	err := Db.SelectOne(&room, "select * from Rooms where slug=$1 and team_id=$2", slug, teamId)
	return &room, err
}
