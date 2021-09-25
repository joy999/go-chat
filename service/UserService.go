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

// 在用户列表中查询用户并返回用户对象
func (this *userService) Find(name string) *model.UserInfo {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if v, ok := this.users[name]; ok {
		return v
	} else {
		return nil
	}
}

// 将某个用户踢下线
func (this *userService) KickOffUser(name string) {
	v := this.Find(name)
	if v != nil {
		this.RemoveUser(name)
	}
}

// 利用用户名及用户的连接对象创建一个用户对象，并将它添加到用户列表中。最后会返回这个新的用户对象
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

// 从用记列表中移除一个用户
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

// 遍历用户列表，会依次遍历用户列表中的每一个用户对象，返回true会继续遍历，false则中止遍历
func (this *userService) Each(fn func(*model.UserInfo) bool) {
	this.lock.RLockFn(func() {
		for _, v := range this.users {
			if !fn(v) {
				return
			}
		}
	})
}
