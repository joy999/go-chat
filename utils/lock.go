package utils

import "sync"

/**
互斥锁封装

增加两个函数式调用方法，方便用于局部加锁
*/

type Locker struct {
	sync.RWMutex
}

func (this *Locker) LockFn(fn func()) {
	this.Lock()
	defer this.Unlock()

	fn()
}

func (this *Locker) RLockFn(fn func()) {
	this.RLock()
	defer this.RUnlock()

	fn()
}
