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

// GM指令，用户状态查询
type Stats struct {
	Username string `json:"username"`
	Long     string `json:"long"`
}

// GM指令，最高频度的词
type Popular struct {
	Words string `json:"words"`
	Cnt   int    `json:"cnt"`
}

// 消息历史记录
type MsgRecord struct {
	Words []string
	Time  int64
}

// 房间列表类
type RoomList struct {
	List []uint `json:"list"`
}
