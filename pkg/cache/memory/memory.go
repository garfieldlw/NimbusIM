package memory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/garfieldlw/NimbusIM/pkg/cache/enum"
	"sync"
	"time"
)

var lock = &sync.Mutex{}
var memoryClient *Memory

type Memory struct {
	cache *bigcache.BigCache
}

func GetMemory() (*Memory, error) {
	if memoryClient != nil {
		return memoryClient, nil
	}

	lock.Lock()
	defer lock.Unlock()

	if memoryClient != nil {
		return memoryClient, nil
	}

	cache, err := bigcache.New(context.Background(), bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted
		LifeWindow: 8 * 60 * time.Minute,

		CleanWindow:        1 * time.Second,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		StatsEnabled:       false,
		Verbose:            true,
		Hasher:             fnv64a{},
		HardMaxCacheSize:   1024,
	})

	if err != nil {
		return nil, err
	}

	memoryClient = &Memory{
		cache: cache,
	}

	return memoryClient, nil
}

func (r *Memory) Get(prefix enum.PrefixEnum, key string) ([]byte, error) {
	key = r.getKey(prefix, key)
	return r.cache.Get(key)
}

func (r *Memory) SetObject(prefix enum.PrefixEnum, key string, object interface{}) error {
	entry, err := json.Marshal(object)
	if err != nil {
		return err
	}

	key = r.getKey(prefix, key)
	return r.cache.Set(key, entry)
}

func (r *Memory) SetString(prefix enum.PrefixEnum, key string, entry string) error {
	key = r.getKey(prefix, key)
	return r.cache.Set(key, []byte(entry))
}

func (r *Memory) Set(prefix enum.PrefixEnum, key string, entry []byte) error {
	key = r.getKey(prefix, key)
	return r.cache.Set(key, entry)
}

func (r *Memory) Del(prefix enum.PrefixEnum, key string) error {
	key = r.getKey(prefix, key)
	return r.cache.Delete(key)
}

func (r *Memory) getKey(prefix enum.PrefixEnum, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}

func IsErrEntryNotFound(err error) bool {
	return errors.Is(
		err,
		bigcache.ErrEntryNotFound,
	)
}
