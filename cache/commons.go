package cache

import (
	"errors"
	"time"

	"github.com/pilillo/igovium/utils"
	"github.com/vmihailenco/msgpack/v5"
)

/*
	Cache is a simple key-value cache with a TTL.
	- L1: Distributed memory cache
	- L2: Local memory cache + DB
	Namespace: a logical grouping of cache entries, i.e., an hashmap - can also include version info, e.g. {NAMESPACE}.{VERSION}
	Key: string identifier - can also include version information of kind {ID}.{VERSION} or {VERSION}.{ID} depending on grouping
	Value: any serializable byte array
*/

// DMCache ... in-memory cache interface
type DMCache interface {
	Init(cfg *utils.DMCacheConfig) error
	Get(key string) (CachePayload, error)
	Put(key string, value CachePayload, ttl time.Duration) error
	Delete(key string) error
	Shutdown() error
}

// DBCache ... db-based cache interface
type DBCache interface {
	Init(cfg *utils.DBCacheConfig) error
	Get(key string) (CachePayload, error)
	Put(key string, value CachePayload, ttl time.Duration) error
	Delete(key string) error
	Size() (int64, error)
}

// CacheEntry ... the entry provided to the CacheService
type CacheEntry struct {
	Key   string       `json:"key"`
	Value CachePayload `json:"value"`
	// generic ttl, same for l1 and l2 (or if not all levels in use)
	TTL *string `json:"ttl,omitempty"`
	// level specific ttl
	TTL1 *string `json:"ttl-l1,omitempty"`
	TTL2 *string `json:"ttl-l2,omitempty"`
}

type CachePayload []byte

// Marshal ... selected marshal method to serialize the cache payload
func (m CachePayload) Marshal() ([]byte, error) {
	return m.MarshalJSON()
}

// Unmarshal ... selected unmarshal method to deserialize the cache payload
func (m CachePayload) Unmarshal(data []byte) error {
	return m.UnmarshalJSON(data)
}

// MarshalJSON returns m as the JSON encoding of m.
func (m CachePayload) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *CachePayload) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

// MarshalBinary using msgpack
func (p *CachePayload) MarshalBinary() ([]byte, error) {
	return msgpack.Marshal(p)
}

// UnmarshalBinary using msgpack
func (p *CachePayload) UnmarshalBinary(data []byte) error {
	return msgpack.Unmarshal(data, p)
}

// CacheService ... a more or less generic interface for a caching service
type CacheService interface {
	Init(cfg *utils.Config) error
	Get(key string) (CachePayload, *utils.Response)
	Put(cacheEntry *CacheEntry) *utils.Response
	Delete(key string) *utils.Response
}
