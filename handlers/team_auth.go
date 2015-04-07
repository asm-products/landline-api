package handlers

import (
	"github.com/asm-products/landline-api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// TeamAuth authenticates the team, either by using the JWT as a token
// (e.g., when interacting with the Landline API on landline.io) or by checking
// the shared secret passed in a Basic Authentication header.
func TeamAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt.ParseFromRequest(
			c.Request,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
		)

		if err != nil {
			sharedSecret, _, ok := c.Request.BasicAuth()

			if !ok {
				c.Fail(401, err)
			}

			team := models.FindTeamBySecret(c.Params.ByName("slug"), sharedSecret)

			c.Set("team", team)
		} else {
			team := models.FindTeamById(token.Claims["id"].(string))

			c.Set("team", team)
		}
	}
}
