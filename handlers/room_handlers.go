package handlers

import (
	"github.com/asm-products/landline-api/models"
	"github.com/gin-gonic/gin"
)

func RoomsIndex(c *gin.Context) {
	result, err := c.Get("user")
	if err != nil {
		c.Fail(500, err)
	}
	user := result.(*models.User)

	var rooms []models.Room
	_, err = models.Db.Select(&rooms, "select * from rooms where team_id=$1", user.TeamId)

	c.JSON(200, gin.H{"rooms": rooms})
}
