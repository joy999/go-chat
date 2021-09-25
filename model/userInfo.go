package model

import "time"

// 用户信息类
type UserInfo struct {
	Conn       *Conn
	Name       string `json:"name"`
	RoomId     uint
	CreateTime time.Time
}
