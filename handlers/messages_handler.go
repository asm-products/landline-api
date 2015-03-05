package handlers

import (
	"github.com/asm-products/landline-api/models"
	"github.com/gin-gonic/gin"
)

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

  messages, err := models.FindMessages(room.Id)
  if err != nil {
    c.Fail(500, err)
  }

  c.JSON(200, gin.H{"messages": messages})
}

func MessagesCreate(c *gin.Context) {
  user, err := GetUserFromContext(c)
  if err != nil {
    c.Fail(500, err)
  }

  room, err := models.FindRoom(c.Params.ByName("room"), user.TeamId)
  if err != nil {
    c.Fail(500, err)
  }

  var json MessageJSON
  c.Bind(&json)

  m, err := SendMessage(user, room, json.Body);
  if (err != nil){
	  c.Fail(500, err)
  }

  c.JSON(200, gin.H{"message": m})
}

func SendMessage(user *models.User, room *models.Room, body string) (*models.Message, error) {
  m := &models.Message{
    RoomId: room.Id,
    UserId: user.Id,
    Body: body,
  }
  Socketio_Server.BroadcastTo(room.Id, "message", m)
  err := models.CreateMessage(m)
  return m, err
}
