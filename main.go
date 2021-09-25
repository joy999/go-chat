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
	ServeRun()
}

// 启动HTTP服务
func ServeRun() {
	http.Handle("/ws", websocket.Handler(accept))      //绑定ws的入口
	http.Handle("/", http.FileServer(http.Dir("web"))) //绑定静态WEB文件目录（客户端）

	err := http.ListenAndServe(":12345", nil) //要改端口就改这里
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

// 接受一个websocket连接
func accept(ws *websocket.Conn) {
	//先获取其姓名
	log.Println("accept")
	var name string

	conn := model.NewConn(ws)

	if err := websocket.Message.Receive(ws, &name); err != nil { //接收用户名，这里没做超时处理，有心人是可以攻击的
		return
	}

	log.Println("name", name)
	user := service.UserService.AddUser(name, conn)
	defer service.UserService.RemoveUser(user) //退出时，清掉用户数据
	log.Println("add user")
	user.Conn.Send("init", nil)
	go router.WSRoute(user) //启动路由的消息循环
	conn.Run()              //启动接发消息循环，同时也会被阻塞在这里，只有当连接中断后才会返回
}
