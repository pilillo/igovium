package cache

import (
	"time"

	"github.com/pilillo/igovium/utils"
)

/*
	Cache is a simple key-value cache with a TTL.
	- L1: Distributed memory cache
	- L2: Local memory cache + DB
	Namespace: a logical grouping of cache entries, i.e., an hashmap - can also include version info, e.g. {NAMESPACE}.{VERSION}
	Key: string identifier - can also include version information of kind {ID}.{VERSION} or {VERSION}.{ID} depending on grouping
	Value: any serializable byte array
*/

type DMCache interface {
	Init(cfg *utils.DMCacheConfig) error
	Get(key string) (interface{}, error)
	Put(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Shutdown() error
}

type DBCache interface {
	Init(cfg *utils.DBCacheConfig) error
	Get(key string) (*DBCacheEntry, error)
	Put(cacheEntry DBCacheEntry) error
	Delete(key string) error
	Size() (int64, error)
}

type CacheEntry struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	TTL   string      `json:"ttl"`
}

type CacheService interface {
	Init(cfg *utils.Config) error
	Get(key string) (interface{}, error)
	Put(cacheEntry *CacheEntry) error
	Delete(key string) error
}
