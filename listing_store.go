package listing_store

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

func Get( key interface{} ) (DirListing, bool) {
	dir,err := DefaultListingStore.store[key]
	return dir,err
}


func Update( key interface{}, listing DirListing ) bool {
	return DefaultListingStore.Update( key, listing )
}

func Statistics() interface{} {
	return DefaultListingStore.Statistics()
}

func (store *ListingMap) Get( key interface{} ) (DirListing, bool) {
	dir,err := store.store[key]
	return dir,err
}


func (store *ListingMap) Update( key interface{}, listing DirListing ) bool {
	store.store[key] = listing
	return true
}

func (store *ListingMap) Statistics() (interface{} ) {
	return struct{
      NumEntries   int
    }{
      NumEntries: len( store.store ),
  }
}
