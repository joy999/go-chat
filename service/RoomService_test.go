package service

import (
	"gochat/model"
	"testing"
)

func TestRoom(t *testing.T) {
	//首先创建一个聊天室供测试使用
	room := RoomService.NewRoom()
	r1 := RoomService.GetRoom(room.Id)
	if r1 != room {
		t.Fatal("Create room failed")
	}

	//这里测试在无用户的状态下的一些错误反馈
	r2 := RoomService.GetRoom(0)
	if r2 != nil {
		t.Fatal("get wrong room failed")
	}

	err := RoomService.AddOneUserToOneRoom("test", room.Id)
	if err.Error() != "User Not Found" {
		t.Fatal("AddOneUserToOneRoom User Not Found test failed")
	}

	err = RoomService.RemoveOneUserFromOneRoom("test", room.Id)
	if err.Error() != "User Not Found" {
		t.Fatal("RemoveOneUserFromOneRoom User Not Found test failed")
	}

	_, err = RoomService.GetRoomList("test")
	if err.Error() != "User Not Found" {
		t.Fatal("GetRoomList User Not Found test failed")
	}

	//模拟创建用户1
	user1 := UserService.AddUser("testuser1", &model.Conn{
		Conn:     nil,
		UserName: "testuser1",
		Sending:  make(chan *model.Cmd, 100),
		Received: make(chan *model.Cmd),
	})

	//测试将用户1加入、移除聊天室，以及接收聊天室列表
	err = RoomService.AddOneUserToOneRoom(user1.Name, 0)
	if err.Error() != "Room Not Found" {
		t.Fatal("AddOneUserToOneRoom Room Not Found test failed")
	}

	err = RoomService.AddOneUserToOneRoom(user1.Name, room.Id)
	if err != nil {
		t.Fatal("AddOneUserToOneRoom failed!", err.Error())
	}
	if sc := <-user1.Conn.Sending; sc.Cmd != "enterroom" {
		t.Error("newroom not send enterroom", sc)
	}

	if rooms, err := RoomService.GetRoomList(user1.Name); err != nil {
		t.Fatal("GetRoomList failed!", err.Error())
	} else {
		if len(rooms) != 1 {
			t.Fatal("Room Number is wrong.", len(rooms))
		}
		if rooms[0].Id != room.Id {
			t.Fatal("GetRoomList failed!", rooms)
		}
	}

	err = RoomService.RemoveOneUserFromOneRoom(user1.Name, 0)
	if err.Error() != "Room Not Found" {
		t.Fatal("RemoveOneUserFromOneRoom Room Not Found test failed")
	}

	err = RoomService.RemoveOneUserFromOneRoom(user1.Name, room.Id)
	if err != nil {
		t.Fatal("RemoveOneUserFromOneRoom failed!", err.Error())
	}
	if sc := <-user1.Conn.Sending; sc.Cmd != "leaveroom" {
		t.Error("newroom not send leaveroom", sc)
	}

	//下面开始聊天室内消息收发测试，这时聊天室里没有用户，所以需要添加用户1与用户2
	user2 := UserService.AddUser("testuser2", &model.Conn{
		Conn:     nil,
		UserName: "testuser2",
		Sending:  make(chan *model.Cmd, 100),
		Received: make(chan *model.Cmd),
	})

	err = RoomService.AddOneUserToOneRoom(user1.Name, room.Id)
	if err != nil {
		t.Fatal("AddOneUserToOneRoom failed!", err.Error())
	}
	if sc := <-user1.Conn.Sending; sc.Cmd != "enterroom" { //将消息消费掉，方便后面测试
		t.Error("newroom not send enterroom", sc)
	}

	err = RoomService.AddOneUserToOneRoom(user2.Name, room.Id)
	if err != nil {
		t.Fatal("AddOneUserToOneRoom failed!", err.Error())
	}
	if sc := <-user2.Conn.Sending; sc.Cmd != "enterroom" { //将消息消费掉，方便后面测试
		t.Error("newroom not send enterroom", sc)
	}

	//开始测试消息发送
	if err = RoomService.SendMsg(0, "", ""); err == nil || err.Error() != "Room Not Found" {
		t.Error("SendMsg to Error Room test failed!", err.Error())
	}

	if err = RoomService.SendMsg(room.Id, user1.Name, "abc"); err != nil {
		t.Error("SendMsg Failed!", err.Error())
	}
	//下面开始测试是否两人都接收到了消息
	if sc := <-user1.Conn.Sending; sc.Cmd != "msg" { //测试用户1是否有收到消息
		t.Error("SendMsg not receive[user1]")
	} else {
		c := sc.Body.(*model.Msg)
		if c.From != user1.Name || c.Msg != "abc" || c.RoomId != room.Id {
			t.Error("SendMsg msg data error!", c)
		}
	}
	if sc := <-user2.Conn.Sending; sc.Cmd != "msg" { //测试用户2是否有收到消息
		t.Error("SendMsg not receive[user1]")
	} else {
		c := sc.Body.(*model.Msg)
		if c.From != user1.Name || c.Msg != "abc" || c.RoomId != room.Id {
			t.Error("SendMsg msg data error!", c)
		}
	}
	//移除用户2后再测试用户2是否可以收到消息
	err = RoomService.RemoveOneUserFromOneRoom(user2.Name, room.Id) //将用户2移除
	if err != nil {
		t.Fatal("RemoveOneUserFromOneRoom failed!", err.Error())
	}
	if sc := <-user2.Conn.Sending; sc.Cmd != "leaveroom" {
		t.Error("newroom not send leaveroom", sc)
	}
	if err = RoomService.SendMsg(room.Id, user1.Name, "abc2"); err != nil {
		t.Error("SendMsg Failed!", err.Error())
	}
	if sc := <-user1.Conn.Sending; sc.Cmd != "msg" { //测试用户1是否有收到消息
		t.Error("SendMsg not receive[user1]")
	} else {
		c := sc.Body.(*model.Msg)
		if c.From != user1.Name || c.Msg != "abc2" || c.RoomId != room.Id {
			t.Error("SendMsg msg data error!", c)
		}
	}
	if len(user2.Conn.Sending) > 0 { //测试用户2是否有收到消息
		t.Error("SendMsg sent to [user2] who is not in list!")
	}

	//测试获取聊天室历史消息
	RoomService.GetAllHistoryMsg(user1, room.Id)
	if sc := <-user1.Conn.Sending; sc.Cmd != "msg" { //测试用户1是否有收到消息
		t.Error("GetAllHistoryMsg failed!")
	} else {
		c := sc.Body.(*model.Msg)
		if c.From != user1.Name || c.Msg != "abc" || c.RoomId != room.Id {
			t.Error("GetAllHistoryMsg failed!", c)
		}
	}
	if sc := <-user1.Conn.Sending; sc.Cmd != "msg" { //测试用户1是否有收到消息
		t.Error("GetAllHistoryMsg failed!")
	} else {
		c := sc.Body.(*model.Msg)
		if c.From != user1.Name || c.Msg != "abc2" || c.RoomId != room.Id {
			t.Error("GetAllHistoryMsg failed!", c)
		}
	}

	err = RoomService.RemoveOneUserFromOneRoom(user1.Name, room.Id) //将用户1移除
	if err != nil {
		t.Fatal("RemoveOneUserFromOneRoom failed!", err.Error())
	}
	UserService.RemoveUser(user1)
	UserService.RemoveUser(user2)
}
