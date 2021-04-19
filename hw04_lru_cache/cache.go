package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       *sync.Mutex
}

type cacheItem struct {
	Key   Key
	Value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mu:       &sync.Mutex{},
	}
}

func (c *lruCache) Set(cacheKey Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, wasInCache := c.items[cacheKey]

	if (c.queue.Len() >= c.capacity) && !wasInCache {
		c.removeLast()
	}

	if wasInCache {
		item.Value = cacheItem{Key: cacheKey, Value: value}
		c.queue.MoveToFront(item)
	} else {
		item = c.queue.PushFront(cacheItem{Key: cacheKey, Value: value})
		c.items[cacheKey] = item
	}
	return wasInCache
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(item)
	return item.Value.(cacheItem).Value, true
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) removeLast() {
	listItem := c.queue.Back()
	item := listItem.Value.(cacheItem)
	delete(c.items, item.Key)
	c.queue.Remove(listItem)
}
