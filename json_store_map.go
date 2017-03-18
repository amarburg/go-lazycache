package lazycache

import (
	"fmt"
	"reflect"
	"sync"
)

//import "net/url"

// import "github.com/amarburg/go-lazyfs"
// import "github.com/amarburg/go-lazyquicktime"

import prom "github.com/prometheus/client_golang/prometheus"

type JSONStore interface {
	Get(key string, ptr interface{}) (bool, error)
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

	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		return fmt.Errorf("MapJSONStore: Update(Pointer " + reflect.TypeOf(value).String() + ")")
	}

	store.store[key] = value
	PromCacheSize.With(prom.Labels{"store": "quicktime"}).Set(float64(len(store.store)))

	//quicktime.DumpTree( DefaultQuicktimeStore.store[ key ].Tree )

	return nil
}

func (store *MapJSONStore) Get(key string, value interface{}) (bool, error) {

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return false, fmt.Errorf("MapJSONStore: Get(Pointer " + rv.String() + ")")
	}

	PromCacheRequests.With(prom.Labels{"store": "quicktime"}).Inc()

	rx, has := store.store[key]

	if has {
		reflect.Indirect(reflect.ValueOf(value)).Set(reflect.ValueOf(rx))
	} else {
		PromCacheMisses.With(prom.Labels{"store": "quicktime"}).Inc()

	}

	return has, nil
}

func CreateMapJSONStore() *MapJSONStore {
	return &MapJSONStore{
		store: make(map[string]interface{}),
	}
}
