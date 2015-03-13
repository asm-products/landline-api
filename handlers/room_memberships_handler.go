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

	membership, err := JoinRoom(user, c.Params.ByName("room"))

	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{"membership": membership})
}

func JoinRoom(user *models.User, roomSlug string) (*models.RoomMembership, error) {
	room, err := models.FindRoom(
		roomSlug,
		user.TeamId,
	)

	if err != nil {
		return nil, err
	}

	r := &models.RoomMembership{
		RoomId: room.Id,
		UserId: user.Id,
	}

	return models.FindOrCreateRoomMembership(r)
}

func RoomMembershipsDelete(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	rid, err := LeaveRoom(user, c.Params.ByName("room"))
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, &gin.H{"deleted": rid})
}

//Removes the user as a member from the room, and returns the room's id.
func LeaveRoom(user *models.User, roomSlug string) (string, error) {
	room, err := models.FindRoom(
		roomSlug,
		user.TeamId,
	)
	if err != nil {
		return "", err
	}

	err = models.DeleteRoomMembership(room.Id, user.Id)
	if err != nil {
		return room.Id, err
	}

	return room.Id, nil
}
