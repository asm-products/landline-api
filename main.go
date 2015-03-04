package main

import (
	"os"
	"runtime"

	"github.com/asm-products/landline-api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/tommy351/gin-cors"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	router := gin.Default()
	router.Use(cors.Middleware(cors.Options{
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "OPTIONS", "POST"},
		AllowOrigins:     []string{"*"},
	}))

	router.OPTIONS("/*cors", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": "ok"})
	})

	// Unauthenticated routes
	router.GET("/sessions/new", handlers.SessionsNew)
	router.GET("/sessions/sso", handlers.SessionsLoginSSO)
	router.POST("/teams", handlers.TeamsCreate)
	router.POST("/teams/:slug", handlers.TeamsLogin)

	// session-keeping for landline.io
	t := router.Group("/teams/:slug")
	t.Use(handlers.TeamAuth(os.Getenv("SECRET")))
	t.GET("/", handlers.TeamsShow)
	t.PUT("/", handlers.TeamsUpdate)

	// authenticated routes
	a := router.Group("/")
	a.Use(handlers.Auth(os.Getenv("SECRET")))
	a.GET("/users", handlers.UsersIndex)
	a.GET("/users/find", handlers.UsersFindOne)
	a.GET("/rooms", handlers.RoomsIndex)
	a.POST("/rooms", handlers.RoomsCreate)
	a.GET("/rooms/:room", handlers.RoomsShow)
	a.GET("/rooms/:room/messages", handlers.MessagesIndex)
	a.POST("/rooms/:room/messages", handlers.MessagesCreate)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	router.Run(":" + port)
}
