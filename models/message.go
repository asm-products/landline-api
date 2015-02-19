package models

import "time"

type Message struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	RoomId    string    `json:"room_id"`
	UserId    string    `json:"user_id"`
	Body      string    `json:"body"`
}
