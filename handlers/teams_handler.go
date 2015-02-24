package handlers

import (
  "encoding/json"

	"github.com/asm-products/landline-api/models"
	"github.com/gin-gonic/gin"
)

type TeamJSON struct {
  Email string `json:"email" binding:"required"`
  EncryptedPassword string `json:"password" binding:"required"`
  SSOSecret string `json:"secret" binding:"required"`
  SSOUrl string `json:"url" binding:"required"`
  Slug string `json:"name" binding:"required"`
}

func TeamsCreate(c *gin.Context) {
  var body TeamJSON
  decoder := json.NewDecoder(c.Request.Body)
  err := decoder.Decode(&body)

	t := &models.Team{
    Email: body.Email,
    EncryptedPassword: body.EncryptedPassword,
    SSOSecret: body.SSOSecret,
    SSOUrl: body.SSOUrl,
    Slug: body.Slug,
  }

  team, err := models.FindOrCreateTeam(t)
  if err != nil {
    panic(err)
  }

  token := GenerateToken(team.Id)

  c.JSON(200, gin.H{"token": token})
}
