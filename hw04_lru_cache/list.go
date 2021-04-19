package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Iterator(fn func(i int, v *ListItem))
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type LinkedList struct {
	frontItem *ListItem
	backItem  *ListItem
	size      int
}

func NewList() List {
	list := new(LinkedList)
	return list
}

func (l LinkedList) Len() int {
	return l.size
}

func (l *LinkedList) Front() *ListItem {
	return l.frontItem
}

func (l *LinkedList) Back() *ListItem {
	return l.backItem
}

func (l *LinkedList) PushFront(v interface{}) *ListItem {
	l.size++
	item := &ListItem{Value: v}

	if l.frontItem != nil {
		l.frontItem.Prev = item
	} else {
		l.backItem = item
	}
	item.Next = l.frontItem
	l.frontItem = item

	return item
}

func (l *LinkedList) PushBack(v interface{}) *ListItem {
	l.size++
	item := &ListItem{Value: v}

	if l.backItem != nil {
		l.backItem.Next = item
	} else {
		l.frontItem = item
	}
	item.Prev = l.backItem
	l.backItem = item

	return item
}

func (l *LinkedList) Remove(i *ListItem) {
	if (l.backItem == nil) && (l.frontItem == nil) {
		return
	}
	l.size--
	front := i.Prev
	back := i.Next
	switch {
	case front != nil && back != nil:
		front.Next = back
		back.Prev = front
	case front == nil && back != nil:
		back.Prev = nil
		l.frontItem = back
	case back == nil && front != nil:
		front.Next = nil
		l.backItem = front
	case front == nil && back == nil:
		l.frontItem = nil
		l.backItem = nil
	}
}

func (l *LinkedList) MoveToFront(i *ListItem) {
	switch {
	case i.Prev == nil:
		// item is first, do nothing
		return
	case i.Next == nil:
		// item is last
		i.Prev.Next = nil
		l.backItem = i.Prev
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.frontItem.Prev = i
	i.Next = l.frontItem
	l.frontItem = i
}

func (l *LinkedList) Iterator(fn func(i int, v *ListItem)) {
	v := l.frontItem
	i := 0
	for v != nil {
		fn(i, v)
		i++
		v = v.Next
	}
}
