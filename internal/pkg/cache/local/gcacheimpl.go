package localcache

import (
	"time"

	"github.com/bluele/gcache"

	"wailik.com/internal/pkg/log"
)

type localCache struct {
	cache gcache.Cache
}

type Options struct{}

func (c *localCache) Set(key interface{}, value interface{}, ttl time.Duration) bool {
	log.Debugf("key:%+v", key)
	log.Debugf("value:%+v", value)
	log.Debugf("ttl:%+v", ttl)

	return c.cache.SetWithExpire(key, value, ttl) == nil
}

func (c *localCache) Get(key interface{}) (interface{}, bool) {
	log.Debugf("key:%+v", key)

	r, err := c.cache.Get(key)
	log.Debugf("value:%+v, error:%+v", r, err)

	return r, (err == nil)
}

func (c *localCache) Del(key interface{}) {
	c.cache.Remove(key)
}

func New(opts Options) (*localCache, error) {
	c := gcache.New(200).Build()
	log.Debugf("cache created")

	return &localCache{cache: c}, nil
}
