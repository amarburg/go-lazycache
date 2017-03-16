package lazycache

import (
	"sync"
)

//import "net/url"

// import "github.com/amarburg/go-lazyfs"
// import "github.com/amarburg/go-lazyquicktime"

import prom "github.com/prometheus/client_golang/prometheus"

type JSONStore interface {
	Get(key string, value interface{}) (bool, error)
	Update(key string, value interface{}) error
	Lock()
	Unlock()
}

type MapJSONStore struct {
	store map[string]interface{}

	mutex sync.Mutex
}

func (store *MapJSONStore) Lock() {
	store.mutex.Lock()
}

func (store *MapJSONStore) Unlock() {
	store.mutex.Unlock()
}

func (store *MapJSONStore) Update(key string, value interface{}) error {

	store.store[key] = value

	PromCacheSize.With(prom.Labels{"store": "quicktime"}).Set(float64(len(store.store)))

	//quicktime.DumpTree( DefaultQuicktimeStore.store[ key ].Tree )

	return nil
}

func (store *MapJSONStore) Get(key string, value interface{}) (bool, error) {
	PromCacheRequests.With(prom.Labels{"store": "quicktime"}).Inc()
	retrieved, has := store.store[key]
  value = &retrieved

	if !has {
		PromCacheMisses.With(prom.Labels{"store": "quicktime"}).Inc()
	}
	return has, nil
}

func CreateMapJSONStore() *MapJSONStore {
	return &MapJSONStore{
		store: make(map[string]interface{}),
	}
}
