package memcache

import (
	"log/slog"
	"time"

	"github.com/dgraph-io/ristretto"
)

type Cache struct {
	cache *ristretto.Cache
}

const cacheDuration = time.Minute * 30

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

func (c *Cache) StoreRSS(key string, rawRSS []byte) {
	c.cache.SetWithTTL(key, rawRSS, int64(len(rawRSS)), cacheDuration)
}

func (c *Cache) GetRSS(key string) ([]byte, bool) {
	val, ok := c.cache.Get(key)
	if !ok {
		return nil, false
	}

	return val.([]byte), true
}

func (c *Cache) StoreHash(key string, hash []byte) {
	c.cache.SetWithTTL("hash:"+key, hash, int64(len(hash)), cacheDuration)
}

func (c *Cache) GetHash(key string) ([]byte, bool) {
	val, ok := c.cache.Get("hash:" + key)
	if !ok {
		return nil, false
	}

	return val.([]byte), true
}

const fileInfoCacheDuration = 7 * 24 * time.Hour

type CachedFileInfo struct {
	ContentType string
	Size        int
}

func (c *Cache) StoreFileInfo(key string, contentType string, size int) {
	c.cache.SetWithTTL("fileinfo:"+key, CachedFileInfo{
		ContentType: contentType,
		Size:        size,
	}, 1, fileInfoCacheDuration)
}

func (c *Cache) GetFileInfo(key string) (contentType string, size int, ok bool) {
	val, found := c.cache.Get("fileinfo:" + key)
	if !found {
		return "", 0, false
	}

	fi := val.(CachedFileInfo)
	return fi.ContentType, fi.Size, true
}
