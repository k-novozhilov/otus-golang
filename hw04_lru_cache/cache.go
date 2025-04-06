package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item := cacheItem{key: key, value: value}

	if listItem, exists := c.items[key]; exists {
		listItem.Value = item
		c.queue.MoveToFront(listItem)
		return true
	}

	if c.queue.Len() == c.capacity {
		back := c.queue.Back()
		if back != nil {
			backItem := back.Value.(cacheItem)
			delete(c.items, backItem.key)
			c.queue.Remove(back)
		}
	}

	listItem := c.queue.PushFront(item)
	c.items[key] = listItem
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if listItem, exists := c.items[key]; exists {
		c.queue.MoveToFront(listItem)
		item := listItem.Value.(cacheItem)
		return item.value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
