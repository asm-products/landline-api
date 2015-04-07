package handlers

import (
	"os"

	"github.com/asm-products/landline-api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func UsersFindOne(c *gin.Context) {
	secret := os.Getenv("SECRET")
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

	c.JSON(200, gin.H{"user": user})
}

func UsersIndex(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	users, err := models.FindRecentlyOnlineUsers(user.TeamId)
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{"users": users})
}

// UsersSearch looks for usernames that resemble the partialUsername
// in the query parameter "u"
func UsersSearch(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.Fail(500, err)
	}

	user, err := GetUserFromContext(c)
	if err != nil {
		c.Fail(500, err)
	}

	partialUsername := c.Request.Form.Get("u")
	users, err := models.SearchUsersByUsernameLike(partialUsername, user.TeamId)
	if err != nil {
		c.Fail(500, err)
	}

	c.JSON(200, gin.H{"users": users})
}

// GetUserFromContext fetches the user that was set in the *gin.Context
// during authorization
func GetUserFromContext(c *gin.Context) (*models.User, error) {
	result, err := c.Get("user")
	if err != nil {
		return nil, err
	}
	user := result.(*models.User)
	return user, nil
}
