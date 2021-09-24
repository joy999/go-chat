package model

import "time"

type UserInfo struct {
	Conn       *Conn
	Name       string `json:"name"`
	RoomId     uint
	CreateTime time.Time
}
