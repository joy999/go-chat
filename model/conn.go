package model

import (
	"encoding/json"
	"gochat/utils"
	"log"
	"time"

	"github.com/gogf/gf/encoding/gjson"
	"golang.org/x/net/websocket"
)

// 连接类，封装了webscoket.Conn等信息
type Conn struct {
	*websocket.Conn
	UserName string

	Sending  chan *Cmd
	Received chan *Cmd

	heartJumpLock  utils.Locker
	heartJumpTimes int //心跳次数
}

//新创一个连接对象
func NewConn(ws *websocket.Conn) *Conn {
	return &Conn{
		Conn:           ws,
		Sending:        make(chan *Cmd, 100),
		Received:       make(chan *Cmd, 0),
		heartJumpTimes: 0,
	}
}

//关闭当前连接
func (this *Conn) Close() {
	if this.Conn != nil {
		this.Conn.Close()
	}
	close(this.Sending)
	close(this.Received)
}

//发送一个消息给当前连接
func (this *Conn) Send(cmd string, body ...interface{}) {
	c := &Cmd{
		Cmd: cmd,
	}
	if len(body) > 0 {
		c.Body = body[0]
	}
	this.Sending <- c
}

//发送系统消息
func (this *Conn) SendSystemMsg(msg string) {
	this.Send("sysmsg", &Msg{
		From: "SYSTEM",
		Msg:  msg,
	})
}

//重置心跳计数
func (this *Conn) ResetHeartJump() {
	this.heartJumpLock.LockFn(func() {
		this.heartJumpTimes = 0
	})
}

//连接的消息循环，这里主要完成消息的收发、心跳下发
func (this *Conn) Run() {
	var (
		c   *Cmd
		ok  bool
		bs  []byte
		err error
	)

	go func() {
		for {
			select {
			case c, ok = <-this.Sending: //写指令
				log.Println("sending...")
				if !ok || c == nil { //故障，退出
					log.Println("Sending Get ERROR", err)
					return
				}
				bs, err = json.Marshal(c)

				if err = websocket.Message.Send(this.Conn, string(bs)); err != nil {
					log.Println("Send ERROR", err)
					this.Conn.Close()
					return
				}

			case <-time.After(time.Second * 15): //心跳
				log.Println("send heartjump...")
				this.Send("heartjump")

				this.heartJumpLock.LockFn(func() {
					this.heartJumpTimes++
					if this.heartJumpTimes >= 3 { //心跳超过3次(180秒)，则关闭连接
						this.Close()
						log.Println("heartJumpTimes >= 3 ERROR", err)
						return
					}
				})
			}
		}
	}()

	for {
		log.Println("wait read...")
		if err = websocket.Message.Receive(this.Conn, &bs); err != nil {
			log.Println("Read ERROR", err)
			this.Conn.Close()
			return
		} else {
			log.Println("Received ", string(bs))
		}
		//解析指令

		if j, err := gjson.DecodeToJson(bs); err == nil {
			//解析无异常，执行
			if err = j.Scan(&c); err != nil {
				log.Println("JSON scan error ", err)
			}

			this.Received <- c
		} else {
			log.Println("JSON parse error ", err)
		}
	}
}
