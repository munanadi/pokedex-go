package pokecache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Add will get the latest cached value
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// fmt.Printf("---Storing key: %s\n", key)

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get will set the cache value
func (c *Cache) Get(key string) ([]byte, bool) {
	// fmt.Printf("---Fetching value for key %s\n", key)

	if _, ok := c.cache[key]; !ok {
		return nil, false
	}
	return c.cache[key].val, true
}

func (c *Cache) reapLoop(timeInterval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.cache {
		timeSinceCreated, err := time.ParseDuration(time.Until(v.createdAt).String())
		if err != nil {
			log.Fatalln("error finding time difference")
		}
		timeSinceCreatedInSeonds := timeSinceCreated.Abs().Seconds()
		if timeSinceCreatedInSeonds > timeInterval.Seconds() {
			// older value, throw it out
			// fmt.Printf("Deleting old value of key %s\n", k)
			delete(c.cache, k)
		}
	}
}

// NewCache will set a cache for a time duration
func NewCache(timeInterval time.Duration) *Cache {
	cache := &Cache{
		cache: map[string]cacheEntry{},
		mu:    sync.Mutex{},
	}

	ticker := time.NewTicker(timeInterval)
	// defer ticker.Stop()

	go func() {
		for range ticker.C {
			// fmt.Printf("\n---Running cleanup function---\n")
			cache.reapLoop(timeInterval)
			fmt.Printf("Pokedex > ")
		}
	}()

	return cache
}
