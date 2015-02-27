package handlers

import (
	"github.com/asm-products/landline-api/models"
	"github.com/gin-gonic/gin"
)

type RoomJSON struct {
	Slug      string    `json:"slug" binding:"required"`
	Topic     string    `json:"topic" binding:"required"`
}

func RoomsIndex(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	var rooms []models.Room
	_, err = models.Db.Select(&rooms, "select * from rooms where team_id=$1", user.TeamId)

	c.JSON(200, gin.H{"rooms": rooms})
}

func RoomsCreate(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	var json RoomJSON
	c.Bind(&json)

	r := &models.Room{
		TeamId: user.TeamId,
		Slug: json.Slug,
		Topic: json.Topic,
	}

	room, err := models.FindOrCreateRoom(r)
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{"room": room})
}

func RoomsShow(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	room, err := models.FindRoom(c.Params.ByName("room"), user.TeamId)
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{"room": room})
}
