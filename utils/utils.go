package utils

import "fmt"

var roomIdLock Locker
var roomId uint

// 生成新的房间ID
func NewRoomId() uint {
	roomIdLock.Lock()
	defer roomIdLock.Unlock()

	roomId++
	return roomId
}

// 将秒数格式化为 00d 00h 00m 00s的格式
func FormatLong(t int64) string {
	s := t % 60
	t /= 60
	m := t % 60
	t /= 60
	h := t % 24
	t /= 24
	d := t

	return fmt.Sprintf("%02dd %02dh %02dm %02ds", d, h, m, s)
}
