package lazycache

import (
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
)

type LocalImageStore struct {
	LocalRoot string
	UrlRoot   string
	logger    kitlog.Logger
	cache     map[string]int

	mutex 		sync.Mutex

	Stats struct {
		cacheRequests int
		cacheMisses   int
	}
}

func (store *LocalImageStore) Has(key string) bool {
	filename := store.LocalRoot + key

	store.mutex.Lock()
	defer store.mutex.Unlock()

	_, has := store.cache[filename]
	if has {
		store.logger.Log("level", "debug", "msg", fmt.Sprintf("Image exists in cache: %s", filename))
		store.cache[filename]++
		return true
	}

	store.logger.Log("level", "debug", "msg", fmt.Sprintf("Checking local image store for \"%s\"", filename))
	_, err := os.Stat(filename)
	if err == nil {
		store.cache[filename] = 1
		return true
	}

	return false
}

func (store *LocalImageStore) Url(key string) (string, bool) {

	store.Stats.cacheRequests++

	if store.Has(key) {
		return store.UrlRoot + key, true
	} else {
		store.Stats.cacheMisses++
		return "", false
	}
}

func RecursiveMkdir(dir string) {
	_, err := os.Stat(dir)
	if err != nil {
		RecursiveMkdir(path.Dir(dir))

		os.Mkdir(dir, 0755)
	}
}

func (store *LocalImageStore) Store(key string, data io.Reader) {
	filename := store.LocalRoot + key
	RecursiveMkdir(path.Dir(filename))

	f, err := os.Create(filename)
	if err != nil {
		store.logger.Log("msg", err.Error(), "type", "error")
	}
	
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.cache[filename] = 1
	io.Copy(f, data)
}

func (store LocalImageStore) Retrieve(key string) (io.Reader, error) {

	f, err := os.Open(store.LocalRoot + key)

	return f, err
}

func (store LocalImageStore) Statistics() interface{} {
	return struct {
		Type          string
		CacheRequests int `json: "cache_requests"`
		CacheMisses   int `json: "cache_misses"`
	}{
		Type:          "local_storage",
		CacheRequests: store.Stats.cacheRequests,
		CacheMisses:   store.Stats.cacheMisses,
	}
}

func (store LocalImageStore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	localPath := path.Join(store.LocalRoot, r.URL.Path)

	if _, err := os.Stat(localPath); err != nil {
		http.Error(w, fmt.Sprintf("Could not find \"%s\"", localPath), 404)
	} else {
		http.ServeFile(w, r, localPath)
	}
}

func CreateLocalStore(localRoot string, addr string) *LocalImageStore {

	// port := 7080
	// addr := fmt.Sprintf("%s:%d", host, port)

	store := &LocalImageStore{
		LocalRoot: localRoot,
		UrlRoot:   addr,
		logger:    kitlog.With(DefaultLogger, "module", "LocalImageStore"),
		cache:     make(map[string]int),
	}

	DefaultLogger.Log("msg",
		fmt.Sprintf("Creating local image store at \"%s\", exposed at \"%s\"\n", store.LocalRoot, store.UrlRoot))

	s := &http.Server{
		Addr:    addr,
		Handler: store,
	}

	go s.ListenAndServe()

	return store
}
