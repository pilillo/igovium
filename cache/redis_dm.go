package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pilillo/igovium/utils"
)

type redisDMCacheType struct {
	client *redis.Client
}

func NewRedisDMCache() DMCache {
	return &redisDMCacheType{}
}

func (c *redisDMCacheType) Init(cfg *utils.DMCacheConfig) error {

	c.client = redis.NewClient(&redis.Options{
		Addr:     cfg.HostAddress,
		Password: cfg.Password,
		DB:       0, //todo: change to cfg.DB?
	})
	var ctx = context.Background()
	if _, err := c.client.Ping(ctx).Result(); err != nil {
		return err
	}
	return nil
}

func (c *redisDMCacheType) Get(key string) (CachePayload, error) {
	var ctx = context.Background()
	value, err := c.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return nil, nil //fmt.Errorf("Key %s does not exist! :: %v", key, err)
	} else if err != nil {
		return nil, err
	} else {
		return CachePayload(value), nil
	}
}

func (c *redisDMCacheType) Put(key string, value CachePayload, ttl time.Duration) error {
	var ctx = context.Background()
	// redis expects string keys and []byte based values (aka strings)
	// marshall whatever go structure we got to a json string, if any
	bVal, err := value.Marshal()
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, bVal, ttl).Err()
}

func (c *redisDMCacheType) Delete(key string) error {
	var ctx = context.Background()
	return c.client.Del(ctx, key).Err()
}

func (c *redisDMCacheType) Shutdown() error {
	var ctx = context.Background()
	//return c.client.ShutdownSave(ctx).Err()
	return c.client.Shutdown(ctx).Err()
}
