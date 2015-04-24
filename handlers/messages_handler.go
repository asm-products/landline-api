package handlers

import (
	"html"
	"time"

	"github.com/asm-products/landline-api/models"
	"github.com/gin-gonic/gin"
)

type BridgeMessageJSON struct {
	Body   string `json:"body" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

type MessageJSON struct {
	Body string `json:"body" binding:"required"`
}

func MessagesIndex(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	room, err := models.FindRoom(c.Params.ByName("room"), user.TeamId)
	if err != nil {
		c.Fail(500, err)
	}

	var messages []models.MessageWithUser
	timestamp := c.Request.URL.Query().Get("t")
	if timestamp != "" {
		time, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			c.Fail(500, err)
		}
		messages, err = models.FindMessagesBeforeTimestamp(room.Id, time)
	} else {
		messages, err = models.FindMessages(room.Id)
	}

	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{"messages": messages})
}

func MessagesBridgeCreate(c *gin.Context) {
	var json BridgeMessageJSON
	c.Bind(&json)

	team, err := GetTeamFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}
	user, err := models.FindUserByExternalIDAndTeam(json.UserID, team.Id)
	if err != nil {
		c.Fail(500, err)
	}

	m, err := SendMessage(
		user,
		c.Params.ByName("room"),
		json.Body,
		"true",
	)
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{"message": m})
}

func MessagesCreate(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	var json MessageJSON
	c.Bind(&json)
	m, err := SendMessage(
		user,
		c.Params.ByName("room"),
		json.Body,
		c.Request.URL.Query().Get("bridge"),
	)
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{"message": m})
}

func SendMessage(user *models.User, roomSlug, body, bridge string) (*models.MessageWithUser, error) {
	room, err := models.FindRoom(roomSlug, user.TeamId)
	if err != nil {
		return nil, err
	}

	m := &models.Message{
		RoomId: room.Id,
		UserId: user.Id,
		Body:   sanitizeBody(body),
	}

	err = models.CreateMessage(m)
	if err != nil {
		return nil, err
	}

	mu := models.NewMessageWithUser(m, user)
	SocketioServer.BroadcastTo(room.Id, "message", mu, roomSlug)
	if bridge != "true" {
		models.PostToTeamWebhook(room.Id, m)
	}
	return mu, err
}

func sanitizeBody(body string) string {
	return html.EscapeString(body)
}

func MessagesHeart(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	h, err := models.CreateMessageHeart(user.Id, c.Params.ByName("message"))
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{"heart": h})
}

func MessagesUnheart(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	err = models.RemoveMessageHeart(user.Id, c.Params.ByName("message"))
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{})
}
