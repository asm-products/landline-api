package main

import (
	"os"
	"runtime"

	"github.com/asm-products/landline-api/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	r := gin.Default()

	// Unauthenticated routes
	r.GET("/sessions/new", handlers.SessionsNew)
	r.GET("/sessions/sso", handlers.SessionsLoginSSO)
	r.POST("/teams", handlers.TeamsCreate)

	// authenticated routes
	a := r.Group("/")
	a.Use(handlers.Auth(os.Getenv("SECRET")))
	a.GET("rooms", handlers.RoomsIndex)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	r.Run(":" + port)
}
