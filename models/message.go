package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"html"
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
	HTMLBody     string    `json:"html_body"`
	Username     string    `db:"username" json:"username"`
	AvatarUrl    string    `db:"avatar_url" json:"avatar_url"`
	LastOnlineAt time.Time `db:"last_online_at" json:"last_online_at"`
	ProfileUrl   string    `db:"profile_url" json:"profile_url"`
}

// NewMessageWithUser parses the message body and joins the user to it
func NewMessageWithUser(message *Message, user *User) *MessageWithUser {
	return &MessageWithUser{
		Id:           message.Id,
		CreatedAt:    message.CreatedAt,
		Body:         message.Body,
		HTMLBody:     ParseMessage(message),
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

	for i := range messages {
		m := &Message{
			Body:   messages[i].Body,
			RoomId: roomId,
		}
		messages[i].HTMLBody = ParseMessage(m)
	}

	return messages, err
}

func CreateMessage(fields *Message) error {
	fields.Body = sanitizeBody(fields.Body)

	if len(fields.Body) == 0 {
		return errors.New("Message body cannot be blank.")
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

// ParseMessage parses the outgoing message according to the following rules:
// - @username
// - #room
func ParseMessage(message *Message) string {
	body := message.Body
	userMentions := utils.ParseUserMentions(body)

	if len(userMentions) > 0 {
		body = replaceUserMentionsWithLinks(message, userMentions)
	}

	roomMentions := utils.ParseRoomMentions(body)

	if len(roomMentions) > 0 {
		body = replaceRoomMentionsWithLinks(message, roomMentions)
	}

	return body
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

func replaceRoomMentionsWithLinks(message *Message, mentions []string) string {
	room := FindRoomById(message.RoomId)
	body := message.Body
	for i := range mentions {
		r, err := FindRoom(mentions[i], room.TeamId)

		if err != nil {
			continue
		}

		link := `<a href="#/rooms/` + r.Slug +
			`" target="_top" title="` + r.Topic + `">#` + r.Slug + `</a>`
		body = strings.Replace(body, `#`+mentions[i], link, 1)
	}

	return body
}

func replaceUserMentionsWithLinks(message *Message, mentions []string) string {
	room := FindRoomById(message.RoomId)
	body := message.Body
	for i := range mentions {
		u, err := FindUserByUsernameAndTeam(mentions[i], room.TeamId)

		if err != nil {
			continue
		}

		link := `<a href="` + u.ProfileUrl + `" target="_top">@` + u.Username + `</a>`
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

func sanitizeBody(body string) string {
	return html.EscapeString(body)
}
