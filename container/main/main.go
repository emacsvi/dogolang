package main

import (
	"container/list"
	"container/ring"
	"fmt"
)

func print(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func main() {
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushFront(0)

	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == 1 {
			l.InsertAfter(1.1, e)
		}
		if e.Value == 2 {
			l.InsertBefore(1.2, e)
		}
	}
	// print(l)
	fmt.Println(l.Front().Value)
	fmt.Println(l.Back().Value)
	l.MoveToBack(l.Front())
	// print(l)

	for e := l.Back(); e != nil; e = e.Prev() {
		fmt.Println(e.Value)
	}

	r := ring.New(3)
	for i := 0; i <= 3; i++ {
		r.Value = i
		r = r.Next()
	}

	for p := r.Next(); p != r; p = p.Next() {
		fmt.Print(p.Value.(int))
		fmt.Print(",")
	}
}
