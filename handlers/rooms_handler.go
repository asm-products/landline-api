package handlers

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/asm-products/landline-api/models"
	"github.com/gin-gonic/gin"
)

type RoomJSON struct {
	Slug  string `json:"slug" binding:"required"`
	Topic string `json:"topic" binding:"required"`
}

func RoomsIndex(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	rooms, err := models.FindRooms(user.TeamId)
	if err != nil {
		c.Fail(500, err)
	}

	memberships, err := models.FindRoomMemberships(user.Id)
	if err != nil {
		c.Fail(500, err)
	}

	unread, err := models.UnreadRooms(user.Id)
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{
		"rooms":        rooms,
		"memberships":  memberships,
		"unread_rooms": unread,
	})
}

func RoomsCreate(c *gin.Context) {
	teamID := models.FindTeamBySlug(c.Params.ByName("slug")).Id
	var json RoomJSON
	c.Bind(&json)

	r := &models.Room{
		TeamId: teamID,
		Slug:   json.Slug,
		Topic:  json.Topic,
	}

	room, err := models.FindOrCreateRoom(r)
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{"room": room})
}

func RoomsDelete(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}
	err = models.DeleteRoom(c.Params.ByName("room"), user.TeamId)
	if err != nil {
		c.Fail(500, err)
	}

	c.String(200, "ok")
}

func RoomsShow(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	roomSlugAndTopic := c.Params.ByName("room")

	r := &models.Room{
		TeamId: user.TeamId,
		Slug:   roomSlugAndTopic,
		Topic:  roomSlugAndTopic,
	}

	room, err := models.FindOrCreateRoom(r)
	if err != nil {
		c.Fail(500, err)
	}

	pixel := createPixel(room.Id, user.Id)

	c.JSON(200, gin.H{"room": room, "pixel": pixel})
}

func RoomsUnread(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	unread, err := models.UnreadRooms(user.Id)
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{"unread": unread})
}

func RoomsUpdate(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	var json RoomJSON
	slug := c.Params.ByName("room")
	c.Bind(&json)

	r := &models.Room{
		Slug:  json.Slug,
		Topic: json.Topic,
	}

	room, err := models.UpdateRoom(slug, user.TeamId, r)
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, room)
}

func createPixel(roomId string, userId string) string {
	url := os.Getenv("RR_URL")
	publicKey := os.Getenv("RR_PUBLIC_KEY")
	sig := createHash(roomId, userId)

	return fmt.Sprintf("%s/t/%s/%s/%s/%s.gif", url, publicKey, roomId, userId, sig)
}

func createHash(roomId string, userId string) string {
	privateKey := os.Getenv("RR_PRIVATE_KEY")
	publicKey := os.Getenv("RR_PUBLIC_KEY")

	h := sha1.New()
	h.Write([]byte(privateKey))
	h.Write([]byte(publicKey))
	h.Write([]byte(roomId))
	h.Write([]byte(userId))

	return hex.EncodeToString(h.Sum(nil))
}
