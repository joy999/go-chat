package main

import (
	"gochat/model"
	"gochat/router"
	"gochat/service"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	http.Handle("/ws", websocket.Handler(accept))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}

func accept(ws *websocket.Conn) {
	//先获取其姓名
	log.Println("accept")
	var name string

	conn := model.NewConn(ws)

	if err := websocket.Message.Receive(ws, &name); err != nil {
		return
	}

	log.Println("name", name)
	user := service.UserService.AddUser(name, conn)
	log.Println("add user")
	go router.WSRoute(user)
	conn.Run()
}
