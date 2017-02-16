package lazycache

//import "net/url"

import "github.com/amarburg/go-lazyfs"
import "github.com/amarburg/go-lazyquicktime"

import prom "github.com/prometheus/client_golang/prometheus"


type QuicktimeMap    map[string]*lazyquicktime.LazyQuicktime

type QuicktimeStore struct {
  store QuicktimeMap
}


var DefaultQuicktimeStore = QuicktimeStore{
  store: make( QuicktimeMap ),
}


func Update( key string, fs lazyfs.FileSource ) (*lazyquicktime.LazyQuicktime,error) {

  var err error
  DefaultQuicktimeStore.store[ key ],err = lazyquicktime.LoadMovMetadata( fs )

  PromCacheSize.With( prom.Labels{"store":"quicktime"}).Set( float64(len(DefaultQuicktimeStore.store)))

  //quicktime.DumpTree( DefaultQuicktimeStore.store[ key ].Tree )

  return  DefaultQuicktimeStore.store[ key ], err
}


func Get( key string ) (*lazyquicktime.LazyQuicktime, bool) {
  PromCacheRequests.With( prom.Labels{"store":"quicktime"}).Inc()
  entry,has := DefaultQuicktimeStore.store[ key ]
  if !has {
    PromCacheMisses.With( prom.Labels{"store":"quicktime"}).Inc()
  }
  return entry, has
}
