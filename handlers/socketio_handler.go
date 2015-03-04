package handlers

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"github.com/asm-products/landline-api/models"
)

var Socketio_Server *socketio.Server

func SetupSocketIOServer(){
	var err error
	Socketio_Server, err = socketio.NewServer(nil)
    if  err != nil  {
        panic ( err )
    }
}

var userConnections map[string]*models.User= make(map[string]*models.User)

func SocketHandler ( c  * gin.Context ) {
    Socketio_Server.On("connection", func(so socketio.Socket) {
		// Since this function is called for every connection,
		// so the user variable is caught in this closure, and is only accessible
		// to this connection.
		var user *models.User

        so.On("auth", func(token string) string {
			fmt.Println("auth: ", token)
			res, err := getUserFromJwt(token, os.Getenv("SECRET"))
			user = res
			if (err != nil){
				return "error: " + err.Error()
			}
			return "success"
        })

		so.On("join", func(roomSlug string) string {
			if (user == nil){
				return "error: not authenticated"
			}
			room, err := models.FindRoom(roomSlug, user.TeamId)
  			if (err != nil) {
    			return "error: " + err.Error()
  			}
			err = so.Join(room.Id)
			if (err != nil) {
				return "error: " + err.Error()
			}
			return "success"
		})

		so.On("leave", func(roomSlug string) string {
			if (user == nil){
				return "error: not authenticated"
			}
			room, err := models.FindRoom(roomSlug, user.TeamId)
  			if (err != nil) {
    			return "error: " + err.Error()
  			}
			err = so.Leave(room.Id)
			if (err != nil) {
				return "error: " + err.Error()
			}
			return "success"
		})

        so.On("disconnection", func() {
            fmt.Println("on disconnect")
        })
    })

    Socketio_Server.On ( "error", func( so socketio.Socket, err error) {
        fmt.Printf ( "[ WebSocket ] Error : %v", err.Error () )
    })

    Socketio_Server.ServeHTTP ( c.Writer, c.Request )
}

