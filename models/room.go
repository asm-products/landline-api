package models

import "time"

type Room struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TeamId    string    `json:"team_id"`
	Slug      string    `json:"slug"`
	Topic     string    `json:"topic"`
}
