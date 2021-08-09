package service

import (
	"log"
	"time"

	"github.com/pilillo/igovium/cache"
	"github.com/pilillo/igovium/commons"
)

type restServiceType struct {
	dmCache cache.DMCache
	dbCache cache.DBCache
}

var restService CacheService = &restServiceType{}

func (s *restServiceType) Init(cfg *commons.Config) error {
	if cfg.DMCacheConfig != nil {
		s.dmCache = cache.NewDMCache()
		err := s.dmCache.Init(cfg.DMCacheConfig)
		if err != nil {
			return err
		}
	}

	if cfg.DBCacheConfig != nil {
		s.dbCache = cache.NewDBCache()
		err := s.dbCache.Init(cfg.DBCacheConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *restServiceType) Get(key string) (interface{}, error) {
	var val interface{}
	var err error

	// lookup 1st level cache (if any)
	if s.dmCache != nil {
		val, err = s.dmCache.Get(key)
		if err != nil {
			return nil, err
		}
		if val != nil {
			return val, nil
		}
	}
	// lookup 2nd level cache (if any)
	if s.dbCache != nil {
		val, err = s.dbCache.Get(key)
		if err != nil {
			return nil, err
		}
		if val != nil {
			return val, nil
		}
	}
	log.Println("Cache miss for key:", key)
	// not found
	return nil, nil
}

func (s *restServiceType) Put(entry *cache.CacheEntry) error {
	var err error

	duration, err := time.ParseDuration(entry.TTL)
	if err != nil {
		return err
	}

	// put on 1st level cache (if any)
	if s.dmCache != nil {
		err = s.dmCache.Put(entry.Key, entry.Value, duration)
		if err != nil {
			return err
		}
	}

	// put on 2nd level cache (if any)
	if s.dbCache != nil {
		cacheEntry := cache.DBCacheEntry{
			Key:   entry.Key,
			Value: entry.Value.([]byte),
			TTL:   duration,
		}
		err = s.dbCache.Put(cacheEntry)
	}

	return nil
}

func (s *restServiceType) Delete(key string) error {
	var err error
	// del on 1st level cache (if any)
	if s.dmCache != nil {
		err = s.dmCache.Delete(key)
		if err != nil {
			return err
		}
	}
	// del on 2nd level cache (if any)
	if s.dbCache != nil {
		err = s.dbCache.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}
