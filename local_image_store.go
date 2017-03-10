package lazycache

import (
	"io"
	"fmt"
	"os"
	kitlog "github.com/go-kit/kit/log"
)


type LocalImageStore struct {
	LocalRoot		string
	UrlRoot     string
	logger      kitlog.Logger

	Stats 		struct{
		cacheRequests			int
		cacheMisses   int
	}
}

func (store LocalImageStore) Has(key string) bool {
	_,err := os.Stat( store.LocalRoot + key )
	return err != nil
}

func (store *LocalImageStore) Url(key string) (string, bool) {

	store.Stats.cacheRequests++

	if store.Has( key ) {
			return store.UrlRoot + key, true
	} else {
		store.Stats.cacheMisses++
		return "", false
	}
}

func (store LocalImageStore) Store(key string, data io.Reader) {
	f, _ := os.Create( store.LocalRoot + key )
	io.Copy(f, data)
}

func (store LocalImageStore) Retrieve(key string) (io.Reader, error) {

	f, err := os.Open( store.LocalRoot + key )

	return f, err
}

func (store LocalImageStore) Statistics() ( interface {} ) {
	return struct{
			Type string
			CacheRequests			int `json: "cache_requests"`
			CacheMisses   int `json: "cache_misses"`
		}{
			Type: "local_storage",
			CacheRequests:  store.Stats.cacheRequests,
			CacheMisses:  store.Stats.cacheMisses,
		}
}

func CreateLocalStore( localRoot string, localUrl string  ) (*LocalImageStore){

  store := &LocalImageStore{
			LocalRoot: localRoot,
			UrlRoot: localUrl,
			logger:  kitlog.With(DefaultLogger, "module", "LocalImageStore"),
	}

  fmt.Printf("Creating local image store at \"%s\", exposed at \"%s\"\n", store.LocalRoot, store.UrlRoot)

	return store
}
