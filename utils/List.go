package utils

type List struct {
	items  []interface{}
	size   int
	header int
	ender  int

	lock Locker
}

func NewList(size int) *List {
	if size == 0 {
		panic("size error")
	}
	o := new(List)
	o.items = make([]interface{}, size)
	o.size = size
	// o.header = 0
	o.ender = 0

	return o
}

func (this *List) Add(s interface{}) {
	this.lock.LockFn(func() {
		this.items[this.ender] = s
		this.ender = (this.ender + 1) % this.size
	})
}

func (this *List) GetAt(pos int) (node interface{}) {
	this.lock.RLockFn(func() {
		node = this.items[pos%this.size]
	})
	return
}
func (this *List) SetAt(pos int, node interface{}) {
	this.lock.LockFn(func() {
		this.items[pos%this.size] = node
	})
}

func (this *List) GetAll() []interface{} {
	s := []interface{}{}
	i := this.ender
	for c := 0; c < this.size; c++ {
		v := this.items[i%this.size]
		if v == nil {
			i++
			continue
		}
		s = append(s, v)
		i++
	}
	return s
}
