package lazycache

import (
	"io"
	"fmt"
	"os"
	"path"
	"net/http"
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
	filename := store.LocalRoot + key
	DefaultLogger.Log("level","debug","msg", fmt.Sprintf("Checking for \"%s\"", filename ) )
	_,err := os.Stat( filename )
	return err == nil
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

func RecursiveMkdir( dir string ){
	_,err := os.Stat( dir )
	if err != nil {
		RecursiveMkdir( path.Dir( dir ) )

		os.Mkdir( dir, 0755 )
	}
}

func (store LocalImageStore) Store(key string, data io.Reader) {
	filename := store.LocalRoot + key
	RecursiveMkdir( path.Dir( filename ))

	f, err := os.Create( filename )
	if err != nil {
		DefaultLogger.Log("msg", err.Error(), "type", "error")
	}

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

func (store LocalImageStore) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	localPath := path.Join( store.LocalRoot, r.URL.Path )

	if _,err := os.Stat( localPath ); err != nil {
		http.Error( w, fmt.Sprintf("Could not find \"%s\"", localPath), 404 )
	} else {
		http.ServeFile( w, r, localPath )
	}
}

func CreateLocalStore( localRoot string, host string  ) (*LocalImageStore){

	port := 7080
	addr := fmt.Sprintf("%s:%d", host, port )

  store := &LocalImageStore{
			LocalRoot: localRoot,
			UrlRoot: "http://" + addr + "/",
			logger:  kitlog.With(DefaultLogger, "module", "LocalImageStore"),
	}

  DefaultLogger.Log("msg",
		fmt.Sprintf("Creating local image store at \"%s\", exposed at \"%s\"\n", store.LocalRoot, store.UrlRoot) )

	s := &http.Server{
		Addr:           addr,
		Handler:        store,
	}

	go s.ListenAndServe()

	return store
}
