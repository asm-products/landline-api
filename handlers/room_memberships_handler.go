package handlers

import (
	"github.com/asm-products/landline-api/models"
	"github.com/gin-gonic/gin"
)

func RoomMembershipsCreate(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	room, err := models.FindRoom(
		c.Params.ByName("room"),
		user.TeamId,
	)
	if err != nil {
		c.Fail(500, err)
	}

	r := &models.RoomMembership{
		RoomId: room.Id,
		UserId: user.Id,
	}

	membership, err := models.FindOrCreateRoomMembership(r)
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{"membership": membership})
}

func RoomMembershipsDelete(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	room, err := models.FindRoom(
		c.Params.ByName("room"),
		user.TeamId,
	)
	if err != nil {
		c.Fail(500, err)
	}

	err = models.DeleteRoomMembership(room.Id, user.Id)

	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, &gin.H{"deleted": room.Id})
}
