package cache

import (
	"time"

	"github.com/pilillo/igovium/commons"
)

type CacheEntry struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	TTL   string      `json:"ttl"`
}

type DMCache interface {
	Init(cfg *commons.DMCacheConfig) error
	Get(key string) (interface{}, error)
	Put(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Shutdown() error
}

type DBCache interface {
	Init(cfg *commons.DBCacheConfig) error
	Get(key string) (*DBCacheEntry, error)
	Put(cacheEntry DBCacheEntry) error
	Delete(key string) error
	Size() (int64, error)
}
