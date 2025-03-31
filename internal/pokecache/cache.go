package pokecache

import(
	"sync"
	"time"
	"fmt"
)

type cacheEntry struct {
	createdAt	time.Time
	Val			[]byte
}

type Cache struct {
	data		map[string]cacheEntry
	dur			time.Duration
	t			*time.Ticker
	mu			*sync.Mutex
}

func NewCache(dur time.Duration) *Cache {
	c := Cache{
		data:	map[string]cacheEntry{},
		dur:	dur,
		t:		time.NewTicker(dur),
		mu:		&sync.Mutex{},
	}

	go c.reapLoop()

	return &c
}

func (c *Cache)Add(key string, val []byte) error {
	c.mu.Lock()
	if _, ok := c.data[key]; ok {
		c.mu.Unlock()
		return fmt.Errorf("Key already in Cache")
	}
	c.mu.Unlock()

	var newEntry cacheEntry
	newEntry.createdAt = time.Now()
	newEntry.Val = val

	c.mu.Lock()
	c.data[key] = newEntry
	c.mu.Unlock()

	return nil
}

func (c *Cache)Get(key string) ([]byte, bool) {
	entry, ok := c.data[key]
	if !ok { return nil, false }
	
	return entry.Val, true
}

func (c *Cache)reapLoop() {
	for {
		<-c.t.C
		c.remOverdue()
	}
}

func (c *Cache)remOverdue() {
	var killTime time.Time

	for key, entry := range c.data {
		killTime = entry.createdAt.Add(c.dur)

		if killTime.After(time.Now()) {
			continue
		}

		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
	}
}
