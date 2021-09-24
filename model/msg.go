package model

//指令结构
type Cmd struct {
	Cmd  string      `json:"cmd"`
	Body interface{} `json:"body"`
}

//业务消息
type Msg struct {
	From   string `json:"from"`
	Msg    string `json:"msg"`
	RoomId uint   `json:"room_id"`
}

type Stats struct {
	Username string `json:"username"`
	Long     string `json:"long"`
}

type Popular struct {
	Words string `json:"words"`
	Cnt   int    `json:"cnt"`
}

type MsgRecord struct {
	Words []string
	Time  int64
}

type RoomList struct {
	List []uint `json:"list"`
}
