package localcache

import (
	"time"

	"github.com/dgraph-io/ristretto"

	"wailik.com/internal/pkg/log"
)

type localCacheR struct {
	cache *ristretto.Cache
}

type OptionsR struct{}

func (c *localCacheR) Set(key interface{}, value interface{}, ttl time.Duration) bool {
	log.Debugf("key:%+v", key)
	log.Debugf("value:%+v", value)

	return c.cache.SetWithTTL(key, value, 1, ttl)
}

func (c *localCacheR) Get(key interface{}) (interface{}, bool) {
	log.Debugf("key:%+v", key)

	r, b := c.cache.Get(key)
	log.Debugf("value:%+v, bool:%+v", r, b)

	return r, b
}

func (c *localCacheR) Del(key interface{}) {
	c.cache.Del(key)
}

func NewR(opts OptionsR) (*localCacheR, error) {
	c, err := ristretto.NewCache(&ristretto.Config{
		MaxCost:     1 << 30,
		NumCounters: 1e7,
		BufferItems: 1,
	})
	if err != nil {
		return nil, err
	}

	return &localCacheR{cache: c}, nil
}
