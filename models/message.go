package models

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/asm-products/landline-api/utils"
	"gopkg.in/gorp.v1"
)

type Message struct {
	Id        string    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	RoomId    string    `db:"room_id" json:"room_id"`
	UserId    string    `db:"user_id" json:"user_id"`
	Body      string    `db:"body" json:"body"`
}

type MessageWithUser struct {
	Id           string    `db:"id" json:"id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	Body         string    `db:"body" json:"body"`
	Username     string    `db:"username" json:"username"`
	AvatarUrl    string    `db:"avatar_url" json:"avatar_url"`
	LastOnlineAt time.Time `db:"last_online_at" json:"last_online_at"`
	ProfileUrl   string    `db:"profile_url" json:"profile_url"`
}

func NewMessageWithUser(message *Message, user *User) *MessageWithUser {
	return &MessageWithUser{
		Id:           message.Id,
		CreatedAt:    message.CreatedAt,
		Body:         message.Body,
		Username:     user.Username,
		AvatarUrl:    user.AvatarUrl,
		LastOnlineAt: user.LastOnlineAt,
		ProfileUrl:   user.ProfileUrl,
	}
}

type UnreadAlert struct {
	Key        string    `json:"key"`
	Recipients *[]string `json:"recipients"`
}

func FindMessages(roomId string) ([]MessageWithUser, error) {
	var messages []MessageWithUser
	_, err := Db.Select(
		&messages,
		`SELECT messages.id, messages.created_at, body, username, avatar_url,
		last_online_at, profile_url
		FROM messages INNER JOIN users ON (users.id = messages.user_id)
		WHERE messages.room_id = $1 ORDER BY messages.created_at ASC
		LIMIT 50`,
		roomId,
	)

	return messages, err
}

func CreateMessage(fields *Message) error {
	mentions := utils.ParseUserMentions(fields.Body)

	if len(mentions) > 0 {
		AlertTeamOfMentions(fields.RoomId, fields.Body, mentions)
		fields.Body = replaceMentionsWithLinks(fields, mentions)
	}

	err := Db.Insert(fields)

	if err != nil {
		return err
	}

	err = registerUnread(fields.RoomId)

	if err != nil {
		return err
	}

	PostToTeamWebhook(fields.RoomId, fields)

	return nil
}

func buildReadraptorRequestBody(roomId string) (*bytes.Reader, error) {
	subscribers, err := Subscribers(roomId)

	if err != nil {
		return nil, err
	}

	alert := UnreadAlert{
		Key:        roomId,
		Recipients: subscribers,
	}

	body, err := json.Marshal(alert)

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), nil
}

func registerUnread(roomId string) error {
	body, err := buildReadraptorRequestBody(roomId)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		os.Getenv("RR_URL")+"/articles",
		body,
	)

	req.SetBasicAuth(os.Getenv("RR_PRIVATE_KEY"), "")

	client := &http.Client{}

	_, err = client.Do(req)

	return err
}

func replaceMentionsWithLinks(fields *Message, mentions []string) string {
	room := FindRoomById(fields.RoomId)
	team := FindTeamById(room.TeamId)
	body := fields.Body
	for i := range mentions {
		u, err := FindUserByUsernameAndTeam(mentions[i], team.Id)

		if err != nil {
			continue
		}

		link := `<a href="` + u.ProfileUrl + `">@` + u.Username + `</a>`
		body = strings.Replace(body, `@`+mentions[i], link, 1)
	}

	return body
}

func (o *Message) PreInsert(s gorp.SqlExecutor) error {
	o.CreatedAt = time.Now()
	o.UpdatedAt = o.CreatedAt

	return nil
}

func (o *Message) PreUpdate(s gorp.SqlExecutor) error {
	o.UpdatedAt = time.Now()
	return nil
}
