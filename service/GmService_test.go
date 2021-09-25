package service

import (
	"gochat/model"
	"testing"
	"time"
)

func TestGmService(t *testing.T) {
	user := UserService.AddUser("testuser", &model.Conn{
		Conn:     nil,
		UserName: "testuser1",
		Sending:  make(chan *model.Cmd, 100),
		Received: make(chan *model.Cmd),
	})
	user.CreateTime = time.Now().Add(-time.Hour * 1)

	GmService.DoCmd(user, "stats testuser")
	if sc := <-user.Conn.Sending; sc.Cmd != "sysmsg" { //测试用户1是否有收到消息
		t.Error("GM: stats testuser failed!")
	} else {
		c := sc.Body.(*model.Msg)
		if c.From != "SYSTEM" || c.Msg != "用户【testuser】在线：00d 01h 00m 00s" {
			t.Error("GM: stats testuser failed!", c)
		}
	}

	time.Sleep(time.Second * 1)
	go func() {
		FilterService.Filter("abc abc abc nnn nnn nn n")
	}()
	time.Sleep(time.Second * 1)
	GmService.DoCmd(user, "popular 2")
	if sc := <-user.Conn.Sending; sc.Cmd != "sysmsg" { //测试用户1是否有收到消息
		t.Error("GM: popular 2 failed!")
	} else {
		c := sc.Body.(*model.Msg)
		if c.From != "SYSTEM" || c.Msg != "2s内出现频率最高的词：abc" {
			t.Error("GM: popular 2 failed!", c)
		}
	}

}
