package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/config"
	"github.com/go-co-op/gocron"
	"github.com/pilillo/igovium/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type dmCacheType struct {
	instance *olric.Olric
	cache    *olric.DMap
}

const (
	// todo: add the concept of namespace to handle multiple maps (multi-tenancy/teams)
	globalMap = "global"
)

func NewDMCache() DMCache {
	return &dmCacheType{}
}

func (c *dmCacheType) Init(cfg *utils.DMCacheConfig) error {
	var err error
	conf := config.New(cfg.Mode)

	ctx, cancel := context.WithCancel(context.Background())
	conf.Started = func() {
		defer cancel()
		log.Println("[INFO] Olric is ready to accept connections")
	}

	c.instance, err = olric.New(conf)
	if err != nil {
		return fmt.Errorf("Failed to create Olric instance: %v", err)
	}

	go func() {
		// Call Start at background. It's a blocker call.
		err = c.instance.Start()
		if err != nil {
			log.Fatalf("olric.Start returned an error: %v", err)
		}
	}()

	<-ctx.Done()

	// create a new cache if it doesn't exist already
	c.cache, err = c.instance.NewDMap(globalMap)
	if err != nil {
		return fmt.Errorf("Failed to create cache DMap: %v", err)
	}
	return nil
}

func (c *dmCacheType) ScheduleStats() error {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Seconds().Do(func() {
		data, err := c.instance.Stats()
		if err != nil {
			log.Printf("[WARN] Failed to get stats: %v", err)
			return
		}
		b, err := json.Marshal(data)
		log.Println(b)
		//data.Partitions[globalMap].DMaps

	})

	s.StartAsync()
	return nil
}

func (c *dmCacheType) Get(key string) (interface{}, error) {
	val, err := c.cache.Get(key)
	if err != nil {
		return nil, fmt.Errorf("Failed to call Get: %v", err)
	}
	if val != nil {
		cachehit.With(prometheus.Labels{"cache": "dm"}).Inc()
	}
	cachemiss.With(prometheus.Labels{"cache": "dm"}).Inc()
	return val, nil
}

func (c *dmCacheType) Put(key string, value interface{}, ttl time.Duration) error {
	err := c.cache.PutEx(key, value, ttl)
	//err := c.cache.Put(key, value)
	if err != nil {
		return fmt.Errorf("Failed to Put on %s: %v", globalMap, err)
	}
	return nil
}

func (c *dmCacheType) Delete(key string) error {
	err := c.cache.Delete(key)
	if err != nil {
		return fmt.Errorf("Failed to Delete %s on %s: %v", key, globalMap, err)
	}
	return nil
}

func (c *dmCacheType) Shutdown() error {
	// leave the cluster
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := c.instance.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("Failed to shutdown Olric: %v", err)
	}
	return nil
}
