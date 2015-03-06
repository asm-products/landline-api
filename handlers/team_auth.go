package handlers

import (
	"github.com/asm-products/landline-api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func TeamAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt.ParseFromRequest(c.Request, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			username, _, ok := c.Request.BasicAuth()

			if !ok {
				c.Fail(401, err)
			}

			team := models.FindTeamBySecret(c.Params.ByName("slug"), username)

			c.Set("team", team)
		} else {
			team := models.FindTeamById(token.Claims["id"].(string))

			c.Set("team", team)
		}
	}
}
