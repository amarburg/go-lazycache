package lazycache

//import "net/url"

import "github.com/amarburg/go-lazyfs"
import "github.com/amarburg/go-lazyquicktime"

import prom "github.com/prometheus/client_golang/prometheus"


type QuicktimeStore interface {
Get( key string ) (*lazyquicktime.LazyQuicktime, bool)
Update( key string, fs lazyfs.FileSource ) (*lazyquicktime.LazyQuicktime,error)
}

type QuicktimeMap    map[string]*lazyquicktime.LazyQuicktime
type MapQuicktimeStore struct {
  store QuicktimeMap
}


var DefaultQuicktimeStore QuicktimeStore


func init() {
  DefaultQuicktimeStore = CreateDefaultQuicktimeStore()
}

func (store *MapQuicktimeStore) Update( key string, fs lazyfs.FileSource ) (*lazyquicktime.LazyQuicktime,error) {

  var err error
  store.store[ key ],err = lazyquicktime.LoadMovMetadata( fs )

  PromCacheSize.With( prom.Labels{"store":"quicktime"}).Set( float64(len(store.store)))

  //quicktime.DumpTree( DefaultQuicktimeStore.store[ key ].Tree )

  return  store.store[ key ], err
}


func (store *MapQuicktimeStore) Get( key string ) (*lazyquicktime.LazyQuicktime, bool) {
  PromCacheRequests.With( prom.Labels{"store":"quicktime"}).Inc()
  entry,has := store.store[ key ]
  if !has {
    PromCacheMisses.With( prom.Labels{"store":"quicktime"}).Inc()
  }
  return entry, has
}

func CreateDefaultQuicktimeStore() (*MapQuicktimeStore) {
    return &MapQuicktimeStore{
      store: make( QuicktimeMap ),
    }
}
