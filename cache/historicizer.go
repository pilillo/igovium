package cache

import (
	"fmt"
	"time"

	"log"

	"github.com/pilillo/igovium/utils"
	"xorm.io/xorm"
)

type dbHistoricizerType struct {
	engine *xorm.Engine
}

func NewDBHistoricizer() *dbHistoricizerType {
	return &dbHistoricizerType{}
}

func (h *dbHistoricizerType) Init(cfg *utils.DBCacheConfig) error {
	var err error
	h.engine, err = xorm.NewEngine(cfg.DriverName, cfg.DataSourceName)
	if err != nil {
		return err
	}
	return nil
}

// GetExpiredAndDelete ... Returns all expired entries and removes them from the database (within the same transaction).
func (h *dbHistoricizerType) GetExpiredAndDelete() ([]DBCacheEntry, error) {

	res, err := h.engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		var expired []DBCacheEntry
		var err error
		now := time.Now().UTC().Unix()
		err = h.engine.
			Table(&DBCacheEntry{}).
			// updated is in secs while ttl is in nanosecs (1 s = 1000000000 ns)
			Where("updated_at + (ttl / 1000000000) <= ?", now).
			//Desc("updated + ttl").
			Find(&expired)
		if err != nil {
			return nil, err
		}

		log.Printf("Found %d expired entries in database", len(expired))

		if len(expired) > 0 {
			var affected int64
			affected, err = h.engine.
				// updated is in secs while ttl is in nanosecs (1 s = 1000000000 ns)
				Where("updated_at + (ttl / 1000000000) <= ?", now).
				Delete(&DBCacheEntry{})
			if err != nil {
				return nil, err
			}
			if affected == 0 || int(affected) != len(expired) {
				return nil, fmt.Errorf("Removed %d entries out of %d identified as expired", affected, len(expired))
			}
		}

		return expired, nil
	})

	if res == nil {
		return nil, err
	}
	return res.([]DBCacheEntry), err
}

func HistoricizeDBCache(config *utils.DBCacheConfig) {
	h := NewDBHistoricizer()
	h.Init(config)
	expired, err := h.GetExpiredAndDelete()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found and removed %d expired entries in database", len(expired))
	// todo: add write to target volume
}
