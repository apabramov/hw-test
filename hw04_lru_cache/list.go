package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	n := &ListItem{Value: v}
	if l.front == nil {
		l.front = n
		l.back = n
	} else {
		n.Next = l.front
		l.front.Prev = n
		l.front = n
	}
	l.len++
	return n
}

func (l *list) PushBack(v interface{}) *ListItem {
	n := &ListItem{Value: v}
	if l.back == nil {
		l.PushFront(n)
	} else {
		n.Prev = l.back
		l.back.Next = n
		l.back = n
	}
	l.len++
	return n
}
func (l *list) Remove(i *ListItem) {
	if i == l.front {
		l.front = l.front.Next
	}
	if i == l.back {
		l.back = l.back.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	i.Prev = nil
	i.Next = nil
	l.len--
}
func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
