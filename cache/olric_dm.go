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

type olricDMCacheType struct {
	instance *olric.Olric
	cache    *olric.DMap
}

const (
	// todo: add the concept of namespace to handle multiple maps (multi-tenancy/teams)
	globalMap = "global"
)

func NewOlricDMCache() DMCache {
	return &olricDMCacheType{}
}

func (c *olricDMCacheType) Init(cfg *utils.DMCacheConfig) error {
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

func (c *olricDMCacheType) ScheduleStats() error {
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

func (c *olricDMCacheType) Get(key string) (CachePayload, error) {
	val, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	if val != nil {
		cachehit.With(prometheus.Labels{"cache": "dm"}).Inc()
	}
	cachemiss.With(prometheus.Labels{"cache": "dm"}).Inc()
	// return cache payload result
	return CachePayload(val.(string)), nil
}

func (c *olricDMCacheType) Put(key string, value CachePayload, ttl time.Duration) error {
	v := string(value)
	err := c.cache.PutEx(key, v, ttl)
	if err != nil {
		return fmt.Errorf("Failed to Put on %s: %v", globalMap, err)
	}
	return nil
}

func (c *olricDMCacheType) Delete(key string) error {
	err := c.cache.Delete(key)
	if err != nil {
		return fmt.Errorf("Failed to Delete %s on %s: %v", key, globalMap, err)
	}
	return nil
}

func (c *olricDMCacheType) Shutdown() error {
	// leave the cluster
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := c.instance.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("Failed to shutdown Olric: %v", err)
	}
	return nil
}
