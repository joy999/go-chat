package api

import (
	"gochat/model"
	"gochat/service"

	"github.com/gogf/gf/encoding/gjson"
)

var RoomApi roomApi

type roomApi struct{}

//获取房间列表
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

//创建房间
func (this roomApi) NewRoom(user *model.UserInfo, data *gjson.Json) {
	room := service.RoomService.NewRoom()
	service.RoomService.AddOneUserToOneRoom(user.Name, room.Id)
	this.GetRoomList(user, data)
}

//切换房间，切换时，会自动退出之前的房间
func (this roomApi) SetRoom(user *model.UserInfo, data *gjson.Json) {
	if user.RoomId > 0 { //先离开原来的房间
		service.RoomService.RemoveOneUserFromOneRoom(user.Name, user.RoomId)
	}
	service.RoomService.AddOneUserToOneRoom(user.Name, data.Var().Uint())
	this.GetHistory(user, data)
}

//离开房间
func (this roomApi) LeaveRoom(user *model.UserInfo, data *gjson.Json) {
	service.RoomService.RemoveOneUserFromOneRoom(user.Name, user.RoomId)
}

//获取所在房间的历史消息，一般不需要使用到此接口，因为进入房间时，历史消息会自动下发
func (this roomApi) GetHistory(user *model.UserInfo, data *gjson.Json) {
	service.RoomService.GetAllHistoryMsg(user, data.Var().Uint())
}

// 发送消息，广播消息给房间内所有人（消息会先被过滤过）
// 若是存在GM指令，则只会将结果返回给查询者，不会广播给房间内所有人
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
