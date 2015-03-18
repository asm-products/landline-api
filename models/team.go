package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/asm-products/landline-api/utils"
	"gopkg.in/gorp.v1"
)

type Team struct {
	Id                string    `db:"id" json:"id"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
	Email             string    `db:"email" json:"email"`
	EncryptedPassword string    `db:"encrypted_password" json:"encrypted_password"`
	SSOSecret         string    `db:"sso_secret" json:"-"`
	SSOUrl            string    `db:"sso_url" json:"sso_url"`
	Slug              string    `db:"slug" json:"slug"`
	WebhookUrl        *string   `db:"webhook_url" json:"webhook_url"`
}

func AlertTeamOfMentions(roomId, body string, mentions []string) error {
	room := FindRoomById(roomId)
	team := FindTeamById(room.TeamId)

	if team.WebhookUrl != nil {
		err := utils.PostMentionsToWebhook(*team.WebhookUrl, team.SSOSecret, body, mentions)

		return err
	}

	return nil
}

func FindOrCreateTeam(fields *Team) (*Team, error) {
	var team Team
	err := Db.SelectOne(&team, "select * from teams where slug=$1", fields.Slug)
	if err == sql.ErrNoRows {
		err = Db.Insert(fields)

		_ = Db.SelectOne(&team, "select * from teams where slug=$1", fields.Slug)

		return fields, err
	}
	return &team, err
}

func FindTeamBySlug(slug string) *Team {
	var team Team
	err := Db.SelectOne(&team, "select * from teams where slug=$1", slug)
	if err != nil {
		panic(err)
	}
	return &team
}

func FindTeamById(id string) *Team {
	var team Team
	err := Db.SelectOne(&team, "select * from Teams where id=$1", id)
	if err != nil {
		panic(err)
	}
	return &team
}

func FindTeamBySecret(slug, secret string) *Team {
	var team Team
	err := Db.SelectOne(
		&team,
		"select * from teams where slug = $1 and sso_secret = $2",
		slug,
		secret,
	)

	if err != nil {
		panic(err)
	}

	return &team
}

func PostToTeamWebhook(roomId string, message *Message) error {
	room := FindRoomById(roomId)
	team := FindTeamById(room.TeamId)
	user, err := FindUser(message.UserId)

	if err != nil {
		return err
	}

	m := utils.Message{
		Body:   message.Body,
		RoomId: roomId,
		UserId: user.ExternalId,
	}

	if team.WebhookUrl != nil {
		err := utils.PostMessageToWebhook(*team.WebhookUrl, team.SSOSecret, m)

		return err
	}

	return nil
}

func UpdateTeam(slug string, fields *Team) (*Team, error) {
	var team Team
	err := Db.SelectOne(&team, "select * from teams where slug=$1", slug)
	if err != nil {
		panic(err)
	}
	team.Email = fields.Email
	team.Slug = fields.Slug
	team.SSOUrl = fields.SSOUrl
	team.SSOSecret = fields.SSOSecret

	_, err = Db.Update(&team)

	return &team, err
}

func ShaString(raw []byte) string {
	hasher := sha256.New()
	hasher.Write(raw)
	return hex.EncodeToString(hasher.Sum(nil))
}

func Sign(secret, payload []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func (o *Team) PreInsert(s gorp.SqlExecutor) error {
	o.CreatedAt = time.Now()
	o.UpdatedAt = o.CreatedAt

	_, _ = FindOrCreateRoom(&Room{
		TeamId: o.Id,
		Slug:   "general",
		Topic:  "general",
	})

	_, _ = FindOrCreateRoom(&Room{
		TeamId: o.Id,
		Slug:   "random",
		Topic:  "random",
	})

	return nil
}

func (o *Team) PreUpdate(s gorp.SqlExecutor) error {
	o.UpdatedAt = time.Now()
	return nil
}
