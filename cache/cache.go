package cache

import (
	"goProj/dataFactory"
	"sync"
)

type Cache struct {
	mutex sync.RWMutex
	data map[string]interface{}
}

func InitCache() *Cache{
	return &Cache{
		data:  make(map[string]interface{}),
	}
}

func (c* Cache) Get(key string) (val interface{}, ok bool) {
	c.mutex.RLock()
	val, ok = c.data[key]
	c.mutex.RUnlock()
	return val, ok
}

func (c* Cache) Store(key string, val interface{}) {
	c.mutex.Lock()
	c.data[key] = val
	c.mutex.Unlock()
}

func (c* Cache) Restore(order *dataFactory.OutputOrder) {
	c.data[order.OrderUID] = order
}
