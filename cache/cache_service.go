package cache

import (
	"log"
	"sync"
	"time"

	_ "sync"

	"github.com/pilillo/igovium/utils"
)

type cacheServiceType struct {
	initialized bool
	initLock    sync.Mutex
	dmCache     DMCache
	dbCache     DBCache
}

var once sync.Once
var cacheService CacheService

// GetCacheService ... Singleton method to instantiate the cache service
func GetCacheService() CacheService {
	once.Do(func() {
		cacheService = &cacheServiceType{initialized: false}
	})
	return cacheService
}

func (s *cacheServiceType) Init(cfg *utils.Config) error {
	s.initLock.Lock()
	defer s.initLock.Unlock()
	if s.initialized {
		return nil
	}

	if cfg.DMCacheConfig != nil {
		s.dmCache = NewDMCache()
		err := s.dmCache.Init(cfg.DMCacheConfig)
		if err != nil {
			return err
		}
	}

	if cfg.DBCacheConfig != nil {
		s.dbCache = NewDBCache()
		err := s.dbCache.Init(cfg.DBCacheConfig)
		if err != nil {
			return err
		}
	}
	// initialization done
	s.initialized = true
	return nil
}

func (s *cacheServiceType) Get(key string) (interface{}, error) {
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

func (s *cacheServiceType) Put(entry *CacheEntry) error {
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
		var byteVal []byte
		byteVal, err = utils.GetBytes(entry.Value)
		cacheEntry := DBCacheEntry{
			Key:   entry.Key,
			Value: byteVal,
			TTL:   duration,
		}
		err = s.dbCache.Put(cacheEntry)
	}

	return nil
}

func (s *cacheServiceType) Delete(key string) error {
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
