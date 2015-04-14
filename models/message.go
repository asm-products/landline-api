package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/asm-products/landline-api/utils"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
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
		WHERE messages.room_id = $1 ORDER BY messages.created_at DESC
		LIMIT 50`,
		roomId,
	)

	return addHTMLBody(roomId, messages), err
}

func FindMessagesBeforeTimestamp(roomId string, timestamp time.Time) ([]MessageWithUser, error) {
	var messages []MessageWithUser
	_, err := Db.Select(
		&messages,
		`SELECT messages.id, messages.created_at, body, username, avatar_url,
		last_online_at, profile_url
		FROM messages INNER JOIN users ON (users.id = messages.user_id)
		WHERE messages.room_id = $1 AND messages.created_at < $2
		ORDER BY messages.created_at DESC
		LIMIT 50`,
		roomId,
		timestamp,
	)

	return addHTMLBody(roomId, messages), err
}

func CreateMessage(fields *Message) error {
	if len(fields.Body) == 0 {
		return errors.New("Message body cannot be blank.")
	}

	err := Db.Insert(fields)

	if err != nil {
		return err
	}

	// ignore Readraptor errors
	_ = registerUnread(fields.RoomId, fields.UserId)

	return nil
}

// ParseMessage parses the outgoing message by catching user and room mentions,
// URLs, and then passing the body through Blackfriday for additional parsing
// and finally Bluemonday for sanitization.
func ParseMessage(message *Message) string {
	body := message.Body

	roomMentions := utils.ParseRoomMentions(body)
	userMentions := utils.ParseUserMentions(body)

	if len(roomMentions) > 0 {
		body = replaceRoomMentionsWithLinks(message, roomMentions)
	}

	if len(userMentions) > 0 {
		body = replaceUserMentionsWithLinks(message, userMentions)
	}

	unsafe := blackfriday.MarkdownCommon([]byte(body))
	safe := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	return string(safe)
}

func addHTMLBody(roomId string, messages []MessageWithUser) []MessageWithUser {
	for i := range messages {
		m := &Message{
			Body:   messages[i].Body,
			RoomId: roomId,
		}

		messages[i].HTMLBody = ParseMessage(m)
	}
	return messages
}

func buildReadraptorRequestBody(roomId, userId string) (*bytes.Reader, error) {
	subscribers, err := SubscribersWithoutUser(roomId, userId)

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

func registerUnread(roomId, userId string) error {
	body, err := buildReadraptorRequestBody(roomId, userId)

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

		link := fmt.Sprintf(
			`<a href="#/rooms/%v" title="%v">#%v</a>`,
			r.Slug,
			r.Topic,
			r.Slug,
		)
		body = strings.Replace(body, `#`+mentions[i], link, 1)
	}

	return body
}

func replaceUrlsWithLinks(message *Message, urls []string) string {
	body := message.Body

	for i := range urls {
		link := fmt.Sprintf(
			`<a href="%v">%v</a>`,
			urls[i],
			urls[i],
		)

		body = strings.Replace(body, urls[i], link, 1)
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

		link := fmt.Sprintf(
			`<a href="%v">@%v</a>`,
			u.ProfileUrl,
			u.Username,
		)

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
