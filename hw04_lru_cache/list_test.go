package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinkedList_Complex(t *testing.T) {
	// внутри тесты не трогал, чтобы было видно что они без изменений
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestLinkedList_PushFront(t *testing.T) {
	t.Run("Push front one", func(t *testing.T) {
		l := NewList()
		v := 55
		l.PushFront(v)
		require.Equal(t, 1, l.Len())
		require.Equal(t, v, l.Front().Value.(int))
		require.Equal(t, v, l.Back().Value.(int))
	})

	t.Run("Push front several", func(t *testing.T) {
		l := NewList()
		itemList := []int{10, 20, 30, 40, 50}
		for _, v := range itemList {
			l.PushFront(v)
		}
		require.Equal(t, len(itemList), l.Len())
		l.Iterator(func(i int, v *ListItem) {
			require.Equal(t, itemList[len(itemList)-i-1], v.Value.(int))
		})
		require.Equal(t, itemList[len(itemList)-1], l.Front().Value.(int))
		require.Equal(t, itemList[0], l.Back().Value.(int))
	})
}

func TestLinkedList_PushBack(t *testing.T) {
	t.Run("Push back one", func(t *testing.T) {
		l := NewList()
		v := 51
		l.PushBack(v)
		require.Equal(t, 1, l.Len())
		require.Equal(t, v, l.Front().Value.(int))
		require.Equal(t, v, l.Back().Value.(int))
	})

	t.Run("Push back several", func(t *testing.T) {
		l := NewList()
		itemList := []int{60, 70, 80, 90, 100}
		for _, v := range itemList {
			l.PushBack(v)
		}
		require.Equal(t, len(itemList), l.Len())
		l.Iterator(func(i int, v *ListItem) {
			require.Equal(t, itemList[i], v.Value.(int))
		})
		require.Equal(t, itemList[0], l.Front().Value.(int))
		require.Equal(t, itemList[len(itemList)-1], l.Back().Value.(int))
	})
}

func TestLinkedList_Remove(t *testing.T) {
	t.Run("Remove single element", func(t *testing.T) {
		l := NewList()
		item := l.PushFront(55)
		l.Remove(item)
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("Remove first", func(t *testing.T) {
		l := NewList()
		itemList := []int{60, 70, 80, 90, 100}
		var rmItem *ListItem
		rmIndex := 0
		for i, v := range itemList {
			item := l.PushBack(v)
			if i == rmIndex {
				rmItem = item
			}
		}
		l.Remove(rmItem)
		itemList = append(itemList[:rmIndex], itemList[rmIndex+1:]...)
		require.Equal(t, len(itemList), l.Len())
		l.Iterator(func(i int, v *ListItem) {
			require.Equal(t, itemList[i], v.Value.(int))
		})
		require.Equal(t, itemList[0], l.Front().Value.(int))
		require.Equal(t, itemList[len(itemList)-1], l.Back().Value.(int))
	})

	t.Run("Remove last", func(t *testing.T) {
		l := NewList()
		itemList := []int{60, 70, 80, 90, 100}
		var rmItem *ListItem
		rmIndex := len(itemList) - 1
		for i, v := range itemList {
			item := l.PushBack(v)
			if i == rmIndex {
				rmItem = item
			}
		}
		l.Remove(rmItem)
		itemList = append(itemList[:rmIndex], itemList[rmIndex+1:]...)
		require.Equal(t, len(itemList), l.Len())
		l.Iterator(func(i int, v *ListItem) {
			require.Equal(t, itemList[i], v.Value.(int))
		})
		require.Equal(t, itemList[0], l.Front().Value.(int))
		require.Equal(t, itemList[len(itemList)-1], l.Back().Value.(int))
	})

	t.Run("Remove middle", func(t *testing.T) {
		l := NewList()
		itemList := []int{60, 70, 80, 90, 100}
		var rmItem *ListItem
		rmIndex := 2
		for i, v := range itemList {
			item := l.PushBack(v)
			if i == rmIndex {
				rmItem = item
			}
		}
		l.Remove(rmItem)
		itemList = append(itemList[:rmIndex], itemList[rmIndex+1:]...)
		require.Equal(t, len(itemList), l.Len())
		l.Iterator(func(i int, v *ListItem) {
			require.Equal(t, itemList[i], v.Value.(int))
		})
		require.Equal(t, itemList[0], l.Front().Value.(int))
		require.Equal(t, itemList[len(itemList)-1], l.Back().Value.(int))
	})
}

func TestLinkedList_MoveToFront(t *testing.T) {
	t.Run("move to front first", func(t *testing.T) {
		l := NewList()
		itemList := []int{60, 70, 80, 90, 100}
		var mvItem *ListItem
		mvIndex := 0
		for i, v := range itemList {
			item := l.PushBack(v)
			if i == mvIndex {
				mvItem = item
			}
		}
		l.MoveToFront(mvItem)
		require.Equal(t, len(itemList), l.Len())
		l.Iterator(func(i int, v *ListItem) {
			require.Equal(t, itemList[i], v.Value.(int))
		})
		require.Equal(t, itemList[0], l.Front().Value.(int))
		require.Equal(t, itemList[len(itemList)-1], l.Back().Value.(int))
	})

	t.Run("move to front last", func(t *testing.T) {
		l := NewList()
		itemList := []int{60, 70, 80, 90, 100}
		var mvItem *ListItem
		mvIndex := len(itemList) - 1
		for i, v := range itemList {
			item := l.PushBack(v)
			if i == mvIndex {
				mvItem = item
			}
		}
		l.MoveToFront(mvItem)
		mvValue := itemList[mvIndex]
		itemList = append(itemList[:mvIndex], itemList[mvIndex+1:]...)
		itemList = append([]int{mvValue}, itemList...)
		require.Equal(t, len(itemList), l.Len())
		l.Iterator(func(i int, v *ListItem) {
			require.Equal(t, itemList[i], v.Value.(int))
		})
		require.Equal(t, itemList[0], l.Front().Value.(int))
		require.Equal(t, itemList[len(itemList)-1], l.Back().Value.(int))
	})

	t.Run("move to front middle", func(t *testing.T) {
		l := NewList()
		itemList := []int{60, 70, 80, 90, 100}
		var mvItem *ListItem
		mvIndex := 2
		for i, v := range itemList {
			item := l.PushBack(v)
			if i == mvIndex {
				mvItem = item
			}
		}
		l.MoveToFront(mvItem)
		mvValue := itemList[mvIndex]
		itemList = append(itemList[:mvIndex], itemList[mvIndex+1:]...)
		itemList = append([]int{mvValue}, itemList...)
		require.Equal(t, len(itemList), l.Len())
		l.Iterator(func(i int, v *ListItem) {
			require.Equal(t, itemList[i], v.Value.(int))
		})
		require.Equal(t, itemList[0], l.Front().Value.(int))
		require.Equal(t, itemList[len(itemList)-1], l.Back().Value.(int))
	})
}
