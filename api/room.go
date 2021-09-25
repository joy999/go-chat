package api

import (
	"gochat/model"
	"gochat/service"

	"github.com/gogf/gf/encoding/gjson"
)

var RoomApi roomApi

type roomApi struct{}

func (this roomApi) GetRoomList(user *model.UserInfo, data *gjson.Json) {
	rs, err := service.RoomService.GetRoomList(user.Name)
	rl := &model.RoomList{
		List: []uint{},
	}
	if err == nil {
		for _, v := range rs {
			rl.List = append(rl.List, v.Id)
		}
		user.Conn.Send("roomlist", rl)
	}
}

func (this roomApi) NewRoom(user *model.UserInfo, data *gjson.Json) {
	room := service.RoomService.NewRoom()
	service.RoomService.AddOneUserToOneRoom(user.Name, room.Id)
	this.GetRoomList(user, data)
}

func (this roomApi) SetRoom(user *model.UserInfo, data *gjson.Json) {
	if user.RoomId > 0 { //先离开原来的房间
		service.RoomService.RemoveOneUserFromOneRoom(user.Name, user.RoomId)
	}
	service.RoomService.AddOneUserToOneRoom(user.Name, data.Var().Uint())
	this.GetHistory(user, data)
}

func (this roomApi) LeaveRoom(user *model.UserInfo, data *gjson.Json) {
	service.RoomService.RemoveOneUserFromOneRoom(user.Name, user.RoomId)
}

func (this roomApi) GetHistory(user *model.UserInfo, data *gjson.Json) {
	service.RoomService.GetAllHistoryMsg(user, data.Var().Uint())
}

func (this roomApi) SendMsg(user *model.UserInfo, data *gjson.Json) {
	msg := data.Var().String()
	if len(msg) == 0 { //空消息不发送
		return
	}
	bs := []byte(msg)
	if bs[0] == '/' {
		service.GmService.DoCmd(user, string(bs[1:]))
		return
	}
	if user.RoomId > 0 {
		service.RoomService.SendMsg(user.RoomId, user.Name, msg)
	}
}
