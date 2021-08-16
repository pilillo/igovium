package cache

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/pilillo/igovium/utils"
	"github.com/prometheus/client_golang/prometheus"

	//_ "github.com/mattn/go-oci8"

	"xorm.io/xorm"
	"xorm.io/xorm/caches"
)

type dbCacheType struct {
	engine *xorm.Engine
}

type DBCacheEntry struct {
	// by default, String is corresponding to varchar(255)
	Key string `xorm:"varchar(255) pk not null unique 'key'"`
	// as time
	//CreatedAt time.Time `xorm:"created"`
	//UpdatedAt time.Time `xorm:"updated"`
	// as local unix timestamps
	CreatedAt int64 `xorm:"created"`
	UpdatedAt int64 `xorm:"updated"`
	Value     []byte
	// shall we need to look/index value, we can deserialize it and add a custom type
	// see https://github.com/go-xorm/tests/blob/master/testCustomTypes.go
	// as delta in nanoseconds
	TTL time.Duration `xorm:"'ttl'"`
}

func NewDBCache() DBCache {
	return &dbCacheType{}
}

func NewInMemoryCache(maxElementSize int) *caches.LRUCacher {
	store := caches.NewMemoryStore()
	cacher := caches.NewLRUCacher(store, maxElementSize)
	return cacher
}

func (c *dbCacheType) Init(cfg *utils.DBCacheConfig) error {
	var err error
	c.engine, err = xorm.NewEngine(cfg.DriverName, cfg.DataSourceName)
	if err != nil {
		return err
	}

	// add local LRU cache if necessary
	if cfg.MaxLocalCacheElementSize > 0 {
		c.engine.SetDefaultCacher(NewInMemoryCache(cfg.MaxLocalCacheElementSize))
	}
	// sync struct with db table schema
	err = c.engine.Sync2(new(DBCacheEntry))
	if err != nil {
		return err
	}

	// start historicize if any conf is set
	if cfg.Historicize != nil {
		err := ScheduleHistoricizeDBCache(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *dbCacheType) Get(key string) (*DBCacheEntry, error) {
	var has bool
	cacheEntry := DBCacheEntry{Key: key}
	has, err := c.engine.Get(&cacheEntry)
	if err != nil {
		return nil, err
	}
	if !has {
		cachemiss.With(prometheus.Labels{"cache": "db"}).Inc()
		return nil, nil
	}
	cachehit.With(prometheus.Labels{"cache": "db"}).Inc()
	return &cacheEntry, nil
}

func (c *dbCacheType) Put(entry DBCacheEntry) error {
	var err error
	var has bool
	// check if key already exists
	lookupEntry := DBCacheEntry{Key: entry.Key}
	has, err = c.engine.Exist(&lookupEntry)
	if err != nil {
		return err
	}
	// if key exists update
	var affected int64
	if has {
		// Update records, bean's non-empty fields are updated contents, condiBean' non-empty filds are conditions CAUTION!
		affected, err = c.engine.Update(&entry, &lookupEntry)
		log.Printf("Updated %d entry", affected)
	} else {
		affected, err = c.engine.Insert(&entry)
		log.Printf("Inserted %d entry", affected)
	}

	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("failed to update cache entry")
	}

	return nil
}

func (c *dbCacheType) Delete(key string) error {
	cacheEntry := DBCacheEntry{Key: key}
	affected, err := c.engine.
		//Where().
		Delete(&cacheEntry)
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("key %s not found, deleted 0 rows", key)
	}
	return nil
}

func (c *dbCacheType) Size() (int64, error) {
	counts, err := c.engine.Count(&DBCacheEntry{})
	if err != nil {
		return -1, err
	}
	return counts, nil
}
