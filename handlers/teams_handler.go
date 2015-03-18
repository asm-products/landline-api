package handlers

import (
	"github.com/asm-products/landline-api/models"
	"github.com/gin-gonic/gin"
)

type TeamJSON struct {
	Email             string  `json:"email" binding:"required"`
	EncryptedPassword string  `json:"password" binding:"required"`
	SSOSecret         string  `json:"secret" binding:"required"`
	SSOUrl            string  `json:"url" binding:"required"`
	Slug              string  `json:"name" binding:"required"`
	WebhookUrl        *string `json:"webhook_url"`
}

type LoginJSON struct {
	EncryptedPassword string `json:"password" binding:"required"`
	Slug              string `json:"name" binding:"required"`
}

func TeamsCreate(c *gin.Context) {
	var json TeamJSON

	c.Bind(&json)

	t := &models.Team{
		Email:             json.Email,
		EncryptedPassword: json.EncryptedPassword,
		SSOSecret:         json.SSOSecret,
		SSOUrl:            json.SSOUrl,
		Slug:              json.Slug,
		WebhookUrl:        json.WebhookUrl,
	}

	team, err := models.FindOrCreateTeam(t)
	if err != nil {
		panic(err)
	}

	token := GenerateToken(team.Id)

	c.JSON(200, gin.H{"token": token})
}

func TeamsShow(c *gin.Context) {
	slug := c.Params.ByName("slug")
	team := models.FindTeamBySlug(slug)

	c.JSON(200, gin.H{
		"email":       team.Email,
		"url":         team.SSOUrl,
		"name":        team.Slug,
		"secret":      team.SSOSecret,
		"webhook_url": team.WebhookUrl,
	})
}

func TeamsLogin(c *gin.Context) {
	var json LoginJSON

	slug := c.Params.ByName("slug")
	team := models.FindTeamBySlug(slug)

	c.Bind(&json)

	if team.EncryptedPassword != json.EncryptedPassword {
		c.String(401, "Unauthorized")
		return
	}

	token := GenerateToken(team.Id)

	c.JSON(200, gin.H{"token": token, "name": team.Slug})
}

func TeamsUpdate(c *gin.Context) {
	var json TeamJSON
	slug := c.Params.ByName("slug")
	c.Bind(&json)

	t := &models.Team{
		Email:     json.Email,
		SSOSecret: json.SSOSecret,
		SSOUrl:    json.SSOUrl,
		Slug:      json.Slug,
	}

	team, err := models.UpdateTeam(slug, t)
	if err != nil {
		panic(err)
	}

	c.JSON(200, team)
}
