package config

import "sync"

type UserCache struct {
	mu    sync.RWMutex
	users map[string]string
}

func NewUserCache() *UserCache {
	return &UserCache{
		users: make(map[string]string),
	}
}

func (c *UserCache) Set(id, name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.users[id] = name
}

func (c *UserCache) Get(id string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	name, exists := c.users[id]
	return name, exists
}

func (c *UserCache) GetAll() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	snapshot := make(map[string]string, len(c.users))
	for k, v := range c.users {
		snapshot[k] = v
	}
	return snapshot
}
