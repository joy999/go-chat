package utils

import "testing"

func TestList(t *testing.T) {
	l := NewList(5)

	l.Add(1)
	all := l.GetAll()
	t.Log(all)
	if len(all) != 1 || all[0] != 1 {
		t.Fatal("1 failed!")
	}

	for i := 2; i <= 5; i++ {
		l.Add(i)
	}
	all = l.GetAll()
	t.Log(all)
	for k, v := range all {
		if k+1 != v.(int) {
			t.Fatal("1 2 3 4 5 failed!")
		}
	}

	l.Add(6)
	all = l.GetAll()
	t.Log(all)
	for i := 2; i <= 6; i++ {
		if i != all[i-2].(int) {
			t.Fatal("2 3 4 5 6 failed")
		}
	}

}
