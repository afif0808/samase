package userservice

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var usersWS map[*websocket.Conn]struct{}
var register chan *websocket.Conn
var unregister chan *websocket.Conn

func GetUserWSs() map[*websocket.Conn]struct{} {

	return usersWS
}

func WebsocketStream() {
	register = make(chan *websocket.Conn)
	unregister = make(chan *websocket.Conn)
	usersWS = map[*websocket.Conn]struct{}{}
	for {
		select {
		case conn := <-register:
			log.Println("wow")
			usersWS[conn] = struct{}{}
			log.Println(usersWS)
		case conn := <-unregister:
			delete(usersWS, conn)
			log.Println(usersWS)
		default:
			time.Sleep(time.Second * 10)

		}
	}
}
