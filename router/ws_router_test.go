package router

import (
	"gochat/model"
	"gochat/service"
	"testing"
)

func TestRouter(t *testing.T) {
	//模拟一个用户，只是这个用户没有连接
	user := service.UserService.AddUser("testuser", &model.Conn{
		Conn:     nil,
		UserName: "testuser",
		Sending:  make(chan *model.Cmd, 100),
		Received: make(chan *model.Cmd),
	})

	go WSRoute(user)

	go func() {
		user.Conn.Received <- &model.Cmd{
			Cmd:  "newroom",
			Body: nil,
		}
	}()

	sc := <-user.Conn.Sending
	if sc.Cmd != "enterroom" {
		t.Error("newroom.enterroom receive error", sc)
	}
	//room_id := sc.Body.(uint)
	sc = <-user.Conn.Sending
	if sc.Cmd != "roomlist" {
		t.Error("newroom.roomlist receive error", sc)
	}

	user.Conn.Close()
}
