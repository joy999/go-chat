package service

import (
	"errors"
	"gochat/model"
	"gochat/utils"
)

var RoomService roomService

type roomService struct {
	rooms map[uint]*model.RoomInfo
	lock  utils.Locker
}

func init() {
	RoomService.init()
}

func (this *roomService) init() {
	this.rooms = map[uint]*model.RoomInfo{}
}

// 新创建一个房间并返回房间对象
func (this *roomService) NewRoom() *model.RoomInfo {
	o := &model.RoomInfo{
		Id:       utils.NewRoomId(),
		UserList: []*model.UserInfo{},
		History:  utils.NewList(50),
	}
	this.lock.LockFn(func() {
		this.rooms[o.Id] = o
	})
	return o
}

// 根据房间ID获取房间对象
func (this *roomService) GetRoom(id uint) (room *model.RoomInfo) {
	this.lock.RLockFn(func() {
		if v, ok := this.rooms[id]; ok {
			room = v
		} else {
			room = nil
		}
	})
	return
}

// 获取房间列表，并将房间列表发送给指定的用户
func (this *roomService) GetRoomList(uname string) (rooms []*model.RoomInfo, err error) {
	user := UserService.Find(uname)
	if user == nil {
		err = errors.New("User Not Found")
		return
	}

	rooms = []*model.RoomInfo{}

	this.lock.RLockFn(func() {
		for _, v := range this.rooms {
			rooms = append(rooms, v)
		}
	})

	return
}

// 将一个用户添加至一个房间中
func (this *roomService) AddOneUserToOneRoom(uname string, room_id uint) error {
	user := UserService.Find(uname)
	if user == nil {
		return errors.New("User Not Found")
	}
	room := this.GetRoom(room_id)
	if room == nil {
		return errors.New("Room Not Found")
	}
	room.Lock.LockFn(func() {
		room.UserList = append(room.UserList, user)
	})
	user.RoomId = room_id

	user.Conn.Send("enterroom", room_id)

	return nil
}

// 将一个用户从一个房间中移除
func (this *roomService) RemoveOneUserFromOneRoom(uname string, room_id uint) error {
	user := UserService.Find(uname)
	if user == nil {
		return errors.New("User Not Found")
	}
	room := this.GetRoom(room_id)
	if room == nil {
		return errors.New("Room Not Found")
	}
	room.Lock.LockFn(func() {
		for k, v := range room.UserList {
			if v == user {
				room.UserList = append(room.UserList[:k], room.UserList[k+1:]...)
				break
			}
		}
	})
	user.RoomId = 0

	user.Conn.Send("leaveroom", nil)

	return nil
}

// 发送消息至指定的房间中
// 注意：这里的发信人(from)只有名字，并不会验证其是否在线
func (this *roomService) SendMsg(room_id uint, from string, msg string) error {
	//脏词过滤
	msg = FilterService.Filter(msg)

	c := &model.Msg{
		From:   from,
		Msg:    msg,
		RoomId: room_id,
	}

	room := this.GetRoom(room_id)

	if room == nil {
		return errors.New("Room Not Found")
	}

	room.History.Add(c)

	room.Lock.RLockFn(func() {
		for _, v := range room.UserList {
			v.Conn.Send("msg", c)
		}
	})

	return nil
}

//获取指定房间中的所有的历史消息（当前最多50条）
func (this *roomService) GetAllHistoryMsg(user *model.UserInfo, room_id uint) {
	room := this.GetRoom(room_id)
	if room == nil {
		return
	}
	msgs := room.History.GetAll()

	for _, c := range msgs {
		user.Conn.Send("msg", c)
	}

}
