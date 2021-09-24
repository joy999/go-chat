package model

import (
	"encoding/json"
	"gochat/utils"
	"time"

	"golang.org/x/net/websocket"
)

type Conn struct {
	*websocket.Conn
	UserName string

	Sending  chan *Cmd
	Received chan *Cmd

	heartJumpLock  utils.Locker
	heartJumpTimes int //心跳次数
}

func NewConn(ws *websocket.Conn) *Conn {
	return &Conn{
		Conn:           ws,
		Sending:        make(chan *Cmd, 100),
		Received:       make(chan *Cmd, 0),
		heartJumpTimes: 0,
	}
}

func (this *Conn) Close() {
	this.Conn.Close()
	close(this.Sending)
}

func (this *Conn) Send(cmd string, body ...interface{}) {
	c := &Cmd{
		Cmd: cmd,
	}
	if len(body) > 0 {
		c.Body = body[0]
	}
	this.Sending <- c
}

func (this *Conn) SendSystemMsg(msg string) {
	this.Send("sysmsg", &Msg{
		From: "SYSTEM",
		Msg:  msg,
	})
}

func (this *Conn) ResetHeartJump() {
	this.heartJumpLock.LockFn(func() {
		this.heartJumpTimes = 0
	})
}

func (this *Conn) Run() {
	var (
		c   *Cmd
		ok  bool
		bs  []byte
		err error
	)
	for {
		select {
		case c, ok = <-this.Sending: //写指令
			if !ok || c == nil { //故障，退出
				return
			}
			bs, err = json.Marshal(c)
			bs = append(bs, '\n')

			if err = websocket.Message.Send(this.Conn, bs); err != nil {
				return
			}

		case <-time.After(time.Second * 15): //心跳
			this.Send("heartjump")

			this.heartJumpLock.LockFn(func() {
				this.heartJumpTimes++
				if this.heartJumpTimes >= 3 { //心跳超过3次(180秒)，则关闭连接
					this.Close()
					return
				}
			})
		default: //读指令
			if err = websocket.Message.Receive(this.Conn, &bs); err != nil {
				return
			}
			//解析指令
			if err = json.Unmarshal(bs, &c); err == nil {
				//解析无异常，执行
				this.Received <- c
			}
		}
	}
}
