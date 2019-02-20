package cache
import "github.com/equalll/mydebug"

type NoneCache struct {
	Cache
}

func NewNoneCache() *NoneCache {mydebug.INFO()
	return &NoneCache{}
}
func (cache *NoneCache) Set(key string, val []byte, life int64) bool {mydebug.INFO()
	return true
}

func (cache *NoneCache) Get(key string) (has bool, data []byte) {mydebug.INFO()
	return false, nil
}
func (cache *NoneCache) Delete(key string) (suc bool) {mydebug.INFO()
	return true
}

func (cache *NoneCache) DeleteAll() (suc bool) {mydebug.INFO()
	return true
}

func (cache *NoneCache) GC() {mydebug.INFO()

}
func (cache *NoneCache) StartGcTimer(sec int64) {mydebug.INFO()
}
