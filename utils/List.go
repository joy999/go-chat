package utils

/*
   这是一个循环队列类，用于需要循环记录的场景
*/
type List struct {
	items  []interface{}
	size   int
	header int
	ender  int

	lock Locker
}

// 生成一个循环队列对象
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

// 增加一个元素到这个队列的写入点
func (this *List) Add(s interface{}) {
	this.lock.LockFn(func() {
		this.items[this.ender] = s
		this.ender = (this.ender + 1) % this.size
	})
}

// 获取指定位置的元素
func (this *List) GetAt(pos int) (node interface{}) {
	this.lock.RLockFn(func() {
		node = this.items[pos%this.size]
	})
	return
}

// 在指定位置存放元素
func (this *List) SetAt(pos int, node interface{}) {
	this.lock.LockFn(func() {
		this.items[pos%this.size] = node
	})
}

//获取所有元素
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
