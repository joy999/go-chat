package utils

import "fmt"

var roomIdLock Locker
var roomId uint

func NewRoomId() uint {
	roomIdLock.Lock()
	defer roomIdLock.Unlock()

	roomId++
	return roomId
}

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
