package model

import "gochat/utils"

type RoomInfo struct {
	Id       uint
	UserList []*UserInfo
	History  *utils.List

	Lock utils.Locker
}
