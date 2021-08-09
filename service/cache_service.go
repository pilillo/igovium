package service

import (
	"github.com/pilillo/igovium/cache"
	"github.com/pilillo/igovium/commons"
)

type CacheService interface {
	Init(cfg *commons.Config) error
	Get(key string) (interface{}, error)
	Put(cacheEntry *cache.CacheEntry) error
	Delete(key string) error
}
