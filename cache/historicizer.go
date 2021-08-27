package cache

import (
	"fmt"
	"os"
	"path"
	"time"

	"log"

	"github.com/go-co-op/gocron"
	"github.com/pilillo/igovium/putters"
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
	// if an engine is set, then we already used it before
	if h.engine != nil {
		return nil
	}
	// otherwise initialize it using the provided conf
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

var h = NewDBHistoricizer()

func HistoricizeDBCache(config *utils.DBCacheConfig) {
	// init xorm db conn - idempotent
	h.Init(config)
	expired, err := h.GetExpiredAndDelete()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found and removed %d expired entries in database", len(expired))
	// todo: add write to target volume
	if len(expired) > 0 {
		// use formatter based on config
		formatManager, err := GetFormatter(config.Historicize.Format)
		if err != nil {
			log.Fatal(err)
		}
		// use date partitioner
		now := time.Now().UTC()
		partName := now.Format(config.Historicize.DatePartitioner)

		// create a tmp folder partition
		partitionPath := path.Join(config.Historicize.TmpDir, partName)
		// create all local partition folders, unless they already exist
		err = os.MkdirAll(partitionPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		// name the file with the current timestamp (does not matter wrt the partition format)
		filename := fmt.Sprintf("%s.%s", fmt.Sprint(now.Unix()), config.Historicize.Format)
		tmpFilePath := path.Join(partitionPath, filename)
		err = formatManager.Save(&expired, tmpFilePath)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Expired entries written to temporary dir %s as %s", tmpFilePath, filename)
		// async put local file to target remote volume if any remote volume is defined
		go putters.Put(config.Historicize.TmpDir, partName, filename, &config.Historicize.RemoteVolumeConfig)
	}
}

var scheduler = gocron.NewScheduler(time.UTC)

func ScheduleHistoricizeDBCache(config *utils.DBCacheConfig) error {
	log.Printf("Running historicize schedule %s", config.Historicize.Schedule)
	// run historicize on schedule as defined in config
	// Schedule:
	// [Minute]			[hour]			[Day_of_the_Month]	[Month_of_the_Year]	[Day_of_the_Week]
	// [0 to 59, or *] 	[0 to 23, or *]	[1 to 31, or *]		[1 to 12, or *]		[0 to 7, with (0 == 7, sunday), or *]
	_, err := scheduler.Cron(config.Historicize.Schedule).Do(HistoricizeDBCache, config)
	if err != nil {
		return err
	}
	// start and continue
	scheduler.StartAsync()
	return nil
}
