package utils

import "sync"

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
