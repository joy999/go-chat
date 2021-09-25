package router

import (
	"gochat/api"
	"gochat/model"

	"github.com/gogf/gf/encoding/gjson"
)

// 这是websocket的路由表，这里用于分发指令至相应的API
func WSRoute(user *model.UserInfo) {
	for {
		cmd, ok := <-user.Conn.Received
		if cmd == nil || !ok {
			return
		}
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
