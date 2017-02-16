package lazycache

import prom "github.com/prometheus/client_golang/prometheus"

type DirListing struct {
	Path        string
	Files       []string
	Directories []string
}

type ListingStore interface {
	Get( interface{} ) (DirListing, bool)
	Update( interface{}, DirListing ) bool
}

type ListingMap struct {
	store    map[interface{}]DirListing
}

var DefaultListingStore = &ListingMap{ store: make( map[interface{}]DirListing ) }

// Convenience wrappers around DefaultListingStore

// func Get( key interface{} ) (DirListing, bool) {
// 	dir,err := DefaultListingStore.store[key]
// 	return dir,err
// }
//
//
// func Update( key interface{}, listing DirListing ) bool {
// 	return DefaultListingStore.Update( key, listing )
// }
//
// func Statistics() interface{} {
// 	return DefaultListingStore.Statistics()
// }

func (store *ListingMap) Get( key interface{} ) (DirListing, bool) {
	dir,has := store.store[key]
	PromCacheRequests.With( prom.Labels{"store":"listing"}).Inc()
	if !has {
		PromCacheMisses.With( prom.Labels{"store":"listing"}).Inc()
	}
	return dir,has
}


func (store *ListingMap) Update( key interface{}, listing DirListing ) bool {
	store.store[key] = listing

	PromCacheSize.With( prom.Labels{"store":"listing"}).Set( float64(len(store.store)) )

	return true
}
