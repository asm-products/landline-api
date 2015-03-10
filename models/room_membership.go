package models

import (
	"database/sql"
	"time"

	"gopkg.in/gorp.v1"
)

type RoomMembership struct {
	Id        string     `db:"id" json:"id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	RoomId    string     `db:"room_id" json:"room_id"`
	UserId    string     `db:"user_id" json:"user_id"`
}

func DeleteRoomMembership(roomId, userId string) error {
	var membership RoomMembership
	err := Db.SelectOne(
		&membership,
		`select * from room_memberships where room_id=$1 and user_id=$2`,
		roomId,
		userId,
	)
	if err != nil {
		panic(err)
	}
	t := time.Now()
	membership.DeletedAt = &t
	_, err = Db.Update(membership)
	return err
}

func FindOrCreateRoomMembership(fields *RoomMembership) (*RoomMembership, error) {
	var membership RoomMembership
	err := Db.SelectOne(
		&membership,
		`select * from room_memberships where room_id=$1 and user_id=$2`,
		fields.RoomId,
		fields.UserId,
	)
	if err == sql.ErrNoRows {
		err = Db.Insert(fields)
		return fields, err
	}
	return &membership, err
}

func FindRoomMemberships(userId string) ([]RoomMembership, error) {
	var memberships []RoomMembership
	_, err := Db.Select(
		&memberships,
		`select room_id from room_memberships where user_id = $1`,
		userId,
	)

	return memberships, err
}

func (o *RoomMembership) PreInsert(s gorp.SqlExecutor) error {
	o.CreatedAt = time.Now()
	o.UpdatedAt = o.CreatedAt
	return nil
}

func (o *RoomMembership) PreUpdate(s gorp.SqlExecutor) error {
	o.UpdatedAt = time.Now()
	return nil
}
