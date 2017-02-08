package quicktime_store

//import "net/url"

import "github.com/amarburg/go-lazyfs"
import "github.com/amarburg/go-lazyquicktime"
//import "github.com/amarburg/go-quicktime"


// type QuicktimeEntry struct {
//   fs      lazyquicktime.LazyQuicktime
// }

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

  //quicktime.DumpTree( DefaultQuicktimeStore.store[ key ].Tree )

  return  DefaultQuicktimeStore.store[ key ], err
}


func Get( key string ) (*lazyquicktime.LazyQuicktime, bool) {
  entry,ok := DefaultQuicktimeStore.store[ key ]
  return entry, ok
}
