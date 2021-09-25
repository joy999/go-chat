package model

import "gochat/utils"

// 房间信息类
type RoomInfo struct {
	Id       uint
	UserList []*UserInfo
	History  *utils.List

	Lock utils.Locker
}
