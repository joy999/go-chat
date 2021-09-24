package service

import (
	"gochat/model"
	"gochat/utils"
)

var ConnService connService

type connService struct {
	connects map[string]*model.Conn
	lock     utils.Locker
}

func init() {
	ConnService.connects = map[string]*model.Conn{}
}

func (this *connService) AddConn(uname string, conn *model.Conn) {
	conn.UserName = uname
	this.lock.LockFn(func() {
		this.connects[uname] = conn
	})
}

func (this *connService) RemoveConn(uname string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.connects, uname)
}
