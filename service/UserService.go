package service

import (
	"gochat/model"
	"gochat/utils"
	"time"
)

//用户管理类，单例
var UserService userService

type userService struct {
	users map[string]*model.UserInfo
	lock  utils.Locker
}

func init() {
	UserService.init()
}

func (this *userService) init() {
	this.users = map[string]*model.UserInfo{}
}

func (this *userService) Find(name string) *model.UserInfo {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if v, ok := this.users[name]; ok {
		return v
	} else {
		return nil
	}
}

func (this *userService) KickOffUser(name string) {
	v := this.Find(name)
	if v != nil {
		this.RemoveUser(name)
	}
}

func (this *userService) AddUser(name string, conn *model.Conn) *model.UserInfo {
	this.KickOffUser(name)

	this.lock.Lock()
	defer this.lock.Unlock()
	info := &model.UserInfo{
		Name:       name,
		Conn:       nil,
		CreateTime: time.Now(),
	}
	this.users[name] = info
	//启动网络数据处理
	ConnService.AddConn(name, conn)
	info.Conn = conn

	return info
}

func (this *userService) RemoveUser(name string) {
	v := this.Find(name)
	if v == nil {
		return
	}
	//从房间中移除
	if v.RoomId > 0 {
		RoomService.RemoveOneUserFromOneRoom(name, v.RoomId)
	}
	//移除用户列表
	this.lock.LockFn(func() {
		delete(this.users, name)
	})
	//移除连接管理
	ConnService.RemoveConn(name)
	//关闭连接
	v.Conn.Close()
}

func (this *userService) Each(fn func(*model.UserInfo) bool) {
	this.lock.RLockFn(func() {
		for _, v := range this.users {
			if !fn(v) {
				return
			}
		}
	})
}
