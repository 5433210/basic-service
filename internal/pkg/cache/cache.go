package cache

import (
	"time"

	localcache "wailik.com/internal/pkg/cache/local"
	remotecache "wailik.com/internal/pkg/cache/remote"
	"wailik.com/internal/pkg/errors"
)

type (
	CacheType     string
	CacheEndpoint string
)

type Cache interface {
	Set(key interface{}, value interface{}, ttl time.Duration) bool
	Get(key interface{}) (interface{}, bool)
	Del(key interface{})
}

type Options struct {
	Type     CacheType
	Endpoint string
}

func New(opts Options) (Cache, error) {
	switch opts.Type {
	case "local":
		c, err := localcache.NewR(localcache.OptionsR{})
		if err != nil {
			return nil, err
		}

		return c, nil

	case "remote":
		c, err := remotecache.New(remotecache.Options{Endpoint: opts.Endpoint})
		if err != nil {
			return nil, err
		}

		return c, nil
	}

	return nil, errors.NewErrorC(errors.ErrCdInvalidCacheType, nil)
}
