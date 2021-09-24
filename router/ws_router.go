package router

import (
	"gochat/api"
	"gochat/model"

	"github.com/gogf/gf/encoding/gjson"
)

func WSRoute(user *model.UserInfo) {
	for {
		cmd := <-user.Conn.Received
		body := gjson.New(cmd.Body)
		switch cmd.Cmd {
		case "newroom":
			api.RoomApi.NewRoom(user, body)
		case "leaveroom":
			api.RoomApi.LeaveRoom(user, body)
		case "setroom":
			api.RoomApi.SetRoom(user, body)
		case "roomlist":
			api.RoomApi.GetRoomList(user, body)
		case "history":
			api.RoomApi.GetHistory(user, body)
		case "heartjump":
			api.SystemApi.HeartJump(user, body)
		case "sendmsg":
			api.RoomApi.SendMsg(user, body)
		}
	}
}
