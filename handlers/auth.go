package handlers

import (
	"github.com/asm-products/landline-api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt.ParseFromRequest(c.Request, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			c.Fail(401, err)
		}

		user, err := models.FindUser(token.Claims["id"].(string))
		if err != nil {
			c.Fail(401, err)
		}

		c.Set("user", user)
	}
}

func getUserFromJwt(webtoken, secret string) (*models.User, error) {
	token, err := jwt.Parse(webtoken, func(token *jwt.Token)(interface{}, error){
		return []byte(secret), nil
	})

	if (err != nil) {
		return nil, err
	}

	user, err := models.FindUser(token.Claims["id"].(string))
	if err != nil {
		return nil, err
	}
	return user, nil
}
