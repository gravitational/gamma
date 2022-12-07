package cache

import (
	"sync"
)

type cache[T any] struct {
	mu sync.RWMutex

	values map[string]T
}

type Cache[T any] interface {
	Get(name string) (T, bool)
	Set(name string, value T)
}

func New[T any]() Cache[T] {
	return &cache[T]{
		values: make(map[string]T),
	}
}

func (c *cache[T]) Set(name string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.values[name] = value
}

func (c *cache[T]) Get(name string) (value T, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok = c.values[name]

	return
}
