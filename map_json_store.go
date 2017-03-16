package lazycache

//import "net/url"

// import "github.com/amarburg/go-lazyfs"
// import "github.com/amarburg/go-lazyquicktime"

import prom "github.com/prometheus/client_golang/prometheus"


// type QuicktimeMap    map[string]*lazyquicktime.LazyQuicktime
// type MapQuicktimeStore struct {
//   store QuicktimeMap
// }



type JSONStore interface {
  Get( key string, value interface{} ) (bool,error)
  Update( key string, value interface{} ) (error)
}

type MapJSONStore struct {
  store map[string]*interface{}
}

func (store *MapJSONStore) MakeKey( key string ) string {
  return key
}


func (store *MapJSONStore) Update( key string, value interface{} ) (error) {

  store.store[ key ]  = &value

  PromCacheSize.With( prom.Labels{"store":"quicktime"}).Set( float64(len(store.store)))

  //quicktime.DumpTree( DefaultQuicktimeStore.store[ key ].Tree )

  return nil
}


func (store *MapJSONStore) Get( key string, value interface{} ) (bool, error) {
  PromCacheRequests.With( prom.Labels{"store":"quicktime"}).Inc()
  value,has := store.store[ key ]
  if !has {
    PromCacheMisses.With( prom.Labels{"store":"quicktime"}).Inc()
  }
  return has, nil
}

func CreateMapJSONStore( ) (*MapJSONStore) {
    return &MapJSONStore{
      store: make( map[string]*interface{} ),
    }
}
