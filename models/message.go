package models

import "time"

type Message struct {
	Id        string    `db:"id" 				 json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	RoomId    string    `db:"room_id"		 json:"room_id"`
	UserId    string    `db:"user_id"	   json:"user_id"`
	Body      string    `db:"body"	 		 json:"body"`
}

type MessageWithUser struct {
	Id        string    `db:"id" 				 json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Body      string    `db:"body"	 		 json:"body"`
	Username	string		`db:"username"	 json:"username"`
	AvatarUrl string		`db:"avatar_url" json:"avatar_url"`
}

func FindMessages(roomId string) ([]MessageWithUser, error) {
	var messages []MessageWithUser
	_, err := Db.Select(
		&messages,
		`SELECT messages.id, messages.created_at, body, username, avatar_url
		FROM messages INNER JOIN users ON (users.id = messages.user_id)
		WHERE messages.room_id = $1 ORDER BY messages.created_at ASC
		LIMIT 50`,
		roomId,
	)

	return messages, err
}

func CreateMessage(fields *Message) error {
	fields.CreatedAt = time.Now()
	fields.UpdatedAt = time.Now()
	err := Db.Insert(fields)

	return err
}
