package cache

import (
	"fmt"
	"log"
	"strings"
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
	var err error
	if cfg.DMCacheConfig != nil {
		if s.dmCache, err = NewDMCacheFromConfig(cfg.DMCacheConfig); err != nil {
			return fmt.Errorf("Error creating DMCache: %v", err)
		}
		if err := s.dmCache.Init(cfg.DMCacheConfig); err != nil {
			return err
		}
	}

	if cfg.DBCacheConfig != nil {
		s.dbCache = NewDBCache()
		if err := s.dbCache.Init(cfg.DBCacheConfig); err != nil {
			return err
		}
	}
	// initialization done
	s.initialized = true
	return nil
}

func (s *cacheServiceType) Get(key string) (CachePayload, *utils.Response) {

	// lookup 1st level cache (if any)
	if s.dmCache != nil {
		dmVal, dmErr := s.dmCache.Get(key)
		if dmErr != nil {
			// both redis and olric return generic errors if key not found so we can't distinguish between them
			if strings.Contains(dmErr.Error(), "not found") {
				return nil, utils.GetNotFoundError(key)
			}
			return nil, utils.GetInternalServerError(dmErr.Error())
		}
		if dmVal != nil {
			log.Printf("Cache hit for %s on DM Cache", key)
			return dmVal, nil
		}
	}
	// lookup 2nd level cache (if any)
	if s.dbCache != nil {
		dbVal, dbErr := s.dbCache.Get(key)
		if dbErr != nil {
			return nil, utils.GetInternalServerError(dbErr.Error())
		}
		if dbVal != nil {
			log.Printf("Cache hit for %s on DB Cache", key)
			return dbVal, nil
		}
	}
	// not found
	return nil, utils.GetNotFoundError(key)
}

func (s *cacheServiceType) Put(entry *CacheEntry) *utils.Response {
	var err error
	var ttl1, ttl2 time.Duration

	// use same ttl if any is set
	if entry.TTL != nil {
		ttl, err := time.ParseDuration(*entry.TTL)
		if err != nil {
			return utils.GetInternalServerError(err.Error())
		}
		ttl1 = ttl
		ttl2 = ttl
	}

	// use level specific ttl if set
	if entry.TTL1 != nil {
		ttl1, err = time.ParseDuration(*entry.TTL1)
		if err != nil {
			return utils.GetInternalServerError(err.Error())
		}
	}
	// use level specific ttl if set
	if entry.TTL2 != nil {
		ttl2, err = time.ParseDuration(*entry.TTL2)
		if err != nil {
			return utils.GetInternalServerError(err.Error())
		}
	}

	// put on 1st level cache (if any)
	if s.dmCache != nil {
		if err = s.dmCache.Put(entry.Key, entry.Value, ttl1); err != nil {
			return utils.GetInternalServerError(err.Error())
		}
	}

	// put on 2nd level cache (if any)
	if s.dbCache != nil {
		if err = s.dbCache.Put(entry.Key, entry.Value, ttl2); err != nil {
			return utils.GetInternalServerError(err.Error())
		}
	}
	// return a successfull response
	return utils.GetPutSuccessfullResponse(entry.Key)
}

func (s *cacheServiceType) Delete(key string) *utils.Response {
	var err error
	// del on 1st level cache (if any)
	if s.dmCache != nil {
		if err = s.dmCache.Delete(key); err != nil {
			return utils.GetInternalServerError(err.Error())
		}
	}
	// del on 2nd level cache (if any)
	if s.dbCache != nil {
		if err = s.dbCache.Delete(key); err != nil {
			return utils.GetInternalServerError(err.Error())
		}
	}
	return utils.GetDeleteSuccessfullResponse(key)
}
