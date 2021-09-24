package service

import "testing"

func TestRoom(t *testing.T) {
	room := RoomService.NewRoom()
	r1 := RoomService.GetRoom(room.Id)
	if r1 != room {
		t.Fatal("Create room failed")
	}

	r2 := RoomService.GetRoom(0)
	if r2 != nil {
		t.Fatal("get wrong room failed")
	}

	err := RoomService.AddOneUserToOneRoom("test", room.Id)
	if err.Error() != "User Not Found" {
		t.Fatal("AddOneUserToOneRoom User Not Found test failed")
	}

	err = RoomService.RemoveOneUserFromOneRoom("test", room.Id)
	if err.Error() != "User Not Found" {
		t.Fatal("RemoveOneUserFromOneRoom User Not Found test failed")
	}
}
