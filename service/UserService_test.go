package service

import (
	"gochat/model"
	"testing"
)

func TestUserService(t *testing.T) {
	user1 := UserService.AddUser("testuser", &model.Conn{
		Conn:     nil,
		UserName: "testuser1",
		Sending:  make(chan *model.Cmd, 100),
		Received: make(chan *model.Cmd),
	})

	if user1 == nil {
		t.Fatal("AddUser Failed!")
	}

	user2 := UserService.AddUser("testuser", &model.Conn{
		Conn:     nil,
		UserName: "testuser1",
		Sending:  make(chan *model.Cmd, 100),
		Received: make(chan *model.Cmd),
	})

	if user2 == nil {
		t.Fatal("AddUser Failed!")
	}

	if user1 == user2 {
		t.Fatal("AddUser Failed!")
	}

	room := RoomService.NewRoom()
	err := RoomService.AddOneUserToOneRoom(user2.Name, room.Id)
	if err != nil || user2.RoomId == 0 {
		t.Fatal("AddOneUserToOneRoom failed!")
	}

	var ok = false
	UserService.Each(func(ui *model.UserInfo) bool {
		if ui.Name == "testuser" {
			ok = true
		}
		return true
	})
	if !ok {
		t.Fatal("Each Failed!")
	}

	UserService.RemoveUser(user2)

	if user2.RoomId > 0 {
		t.Fatal("RemoveUser failed!")
	}

	UserService.Each(func(ui *model.UserInfo) bool {
		if ui.Name == "testuser" {
			t.Fatal("testuser exits! It's wrong")
		}
		return true
	})

	//这个时候user2已经不在用户管理列表中，这里应该返回nil
	if u := UserService.getUser(user2); u != nil {
		t.Fatal("getUser failed")
	}
	//传入非string及*model.UserInfo也应该返回null
	if u := UserService.getUser(2); u != nil {
		t.Fatal("getUser failed")
	}
}
