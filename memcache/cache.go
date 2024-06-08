package memcache

import (
	"log/slog"
	"time"

	"github.com/dgraph-io/ristretto"
)

type Cache struct {
	cache *ristretto.Cache
}

func NewCache() *Cache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,
		MaxCost:     100_000_000, // 100 MB
		BufferItems: 64,
		Metrics:     true,
	})
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(time.Hour)
			slog.Info("cache metrics",
				"cost-added", cache.Metrics.CostAdded(),
				"cost-evicted", cache.Metrics.CostEvicted(),
			)
		}
	}()

	return &Cache{
		cache: cache,
	}
}

func (c *Cache) StoreRSS(id int, rawRSS []byte) {
	c.cache.SetWithTTL(id, rawRSS, int64(len(rawRSS)), time.Minute*15)
}

func (c *Cache) GetRSS(id int) ([]byte, bool) {
	val, ok := c.cache.Get(id)
	if !ok {
		return nil, false
	}

	return val.([]byte), true
}
