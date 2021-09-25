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
	info.Conn = conn

	return info
}

//这个方法在获取用户对象的同时，也可以校验用户对象是否是被管理的状态
func (this *userService) getUser(user interface{}) *model.UserInfo {
	switch v := user.(type) {
	case string:
		return this.Find(v)
	case *model.UserInfo:
		v0 := this.Find(v.Name)
		if v0 == v {
			return v
		} else {
			return nil
		}
	default:
		return nil
	}
}

func (this *userService) RemoveUser(user interface{}) {
	v := this.getUser(user)
	if v == nil {
		return
	}
	name := v.Name
	//从房间中移除
	if v.RoomId > 0 {
		RoomService.RemoveOneUserFromOneRoom(name, v.RoomId)
	}
	//移除用户列表
	this.lock.LockFn(func() {
		delete(this.users, name)
	})
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
