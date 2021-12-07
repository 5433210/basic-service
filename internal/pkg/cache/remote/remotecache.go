package remotecache

import (
	"time"
)

type remoteCache struct{}

type Options struct {
	Endpoint string
}

func New(opts Options) (*remoteCache, error) {
	return &remoteCache{}, nil
}

func (c *remoteCache) Set(key interface{}, value interface{}, ttl time.Duration) bool {
	return false
}

func (c *remoteCache) Get(key interface{}) (interface{}, bool) {
	return nil, false
}

func (c *remoteCache) Del(key interface{}) {
}
