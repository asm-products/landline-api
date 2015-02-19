package models

import "time"

type User struct {
	Id         string    `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	TeamId     string    `json:"team_id"`
	Email      string    `json:"email"`
	ProfileUrl string    `json:"profile_url"`
	RealName   string    `json:"real_name"`
	Username   string    `json:"username"`
}
