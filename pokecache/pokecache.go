package pokecache

import (
	"sync"
	"time"
)

type Cache struct{
	items map[string]CacheEntry
	sc sync.RWMutex
	interval time.Duration
}

type CacheEntry struct{
	createdAt time.Time
	val []byte
}

func NewCache(interval time.Duration)*Cache{
	c := Cache{
		interval: interval,
		sc: sync.RWMutex{},
		items: make(map[string]CacheEntry),
	}
	go c.reapLoop()
	return &c
}

func (cache *Cache)Add(key string, val []byte){
	cache.sc.Lock()
	defer cache.sc.Unlock()
	cache.items[key] = CacheEntry{
		createdAt: time.Now(),
		val: val,
	}
}

func (cache *Cache)Get(key string)([]byte, bool){
	cache.sc.RLock()
	defer cache.sc.RUnlock()
	if val, ok := cache.items[key]; ok{
		return val.val, ok
	}
	return nil, false
}

func (cache *Cache) reapLoop(){
	ticker := time.Tick(cache.interval)
	for range ticker{
		cache.sc.Lock()
		for k,v := range cache.items{
			if v.createdAt.Add(time.Second * 30).Compare(time.Now()) <= 0{
				delete(cache.items, k)
			}
		}
		cache.sc.Unlock()
	}
}
