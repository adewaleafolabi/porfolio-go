package pricing

import (
	"sync"
	"time"
)

type priceCacheItem struct {
	value      float64
	lastAccess int64
}

type PriceCache struct {
	items map[string]*priceCacheItem
	mutex sync.Mutex
}

func NewPriceCache(ttl int) *PriceCache {
	p := &PriceCache{
		items: make(map[string]*priceCacheItem),
	}

	go func() {
		for now := range time.Tick(time.Second) {
			p.mutex.Lock()
			for k, v := range p.items {
				if now.Unix()-v.lastAccess > int64(ttl) {
					delete(p.items, k)
				}
			}
			p.mutex.Unlock()
		}
	}()

	return p
}

func (p *PriceCache) Get(key string) float64 {
	price := 0.0

	p.mutex.Lock()

	if item, ok := p.items[key]; ok {
		price = item.value
		item.lastAccess = time.Now().Unix()
	}

	p.mutex.Unlock()

	return price
}

func (p *PriceCache) Set(key string, value float64) {
	p.mutex.Lock()
	p.items[key] = &priceCacheItem{
		value:      value,
		lastAccess: time.Now().Unix(),
	}
	p.mutex.Unlock()
}
