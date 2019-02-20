package cache
import "github.com/equalll/mydebug"

type Cache interface {
	Set(key string, val []byte, life int64) (suc bool)
	Get(key string) (has bool, data []byte)
	Delete(key string) (suc bool)
	DeleteAll() (suc bool)
	GC()
	StartGcTimer(sec int64)
}

type Data struct {
	Key        string
	Data       []byte
	CreateTime int64
	Life       int64
}

var defaultCache Cache = new(NoneCache)

func SetDefaultCacheHandler(cache Cache) {mydebug.INFO()
	defaultCache = cache
}

func Set(key string, val []byte, life int64) (suc bool) {mydebug.INFO()
	return defaultCache.Set(key, val, life)
}

func Get(key string) (has bool, data []byte) {mydebug.INFO()
	return defaultCache.Get(key)
}

func Delete(key string) (suc bool) {mydebug.INFO()
	return defaultCache.Delete(key)
}
func GC() {mydebug.INFO()
	defaultCache.GC()
}
