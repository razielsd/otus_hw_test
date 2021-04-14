package hw04lrucache

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic - only inserts", func(t *testing.T) {
		c := NewCache(3)
		for i := 1; i <= 4; i++ {
			c.Set(Key(fmt.Sprintf("a%d", i)), 100+i)
		}
		_, ok := c.Get("a1")
		require.False(t, ok, "item must be purged")

		for i := 2; i <= 4; i++ {
			val, ok := c.Get(Key(fmt.Sprintf("a%d", i)))
			require.True(t, ok)
			require.Equal(t, 100+i, val)
		}
	})

	t.Run("purge logic - check access", func(t *testing.T) {
		c := NewCache(3)
		for i := 1; i <= 3; i++ {
			c.Set(Key(fmt.Sprintf("a%d", i)), 100+i)
		}
		_, ok := c.Get("a1")
		require.True(t, ok, "first item must be saved")

		c.Set("a4", 104)

		_, ok = c.Get("a2")
		require.False(t, ok, "second item (key=a2) must be purged after add 4th item")

		keyList := map[Key]int{"a1": 101, "a3": 103, "a4": 104}
		for k, v := range keyList {
			cv, ok := c.Get(k)
			require.True(t, ok, "Purged item with key: "+k+", but must be saved")
			require.Equal(t, v, cv, "Wrong value from cache")
		}
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
