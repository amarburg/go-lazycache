package lazycache

//import "net/url"

import (
	"encoding/json"
	"fmt"
	"sync"
)

import prom "github.com/prometheus/client_golang/prometheus"

import "github.com/mediocregopher/radix.v2/pool"

type RedisJsonStore struct {
	prefix string
	pool   *pool.Pool
	mutex  sync.Mutex
}

func (red *RedisJsonStore) makeKey(key string) string {
	return red.prefix + "/" + key
}

func (red *RedisJsonStore) Lock() {
	red.mutex.Lock()
}

func (red *RedisJsonStore) Unlock() {
	red.mutex.Unlock()
}

func (red *RedisJsonStore) Update(key string, value interface{}) error {
	var err error

	// metadata,err := lazyquicktime.LoadMovMetadata( fs )
	// if err != nil {
	//   DefaultLogger.Log("level", "error", "msg", fmt.Sprintf("Error parsing Quicktime metadata for %s: %s", key, err.Error() ) )
	// return nil, err
	// }

	// PromCacheSize.With( prom.Labels{"store":"quicktime"}).Set( float64(len(red.store)))

	conn, err := red.pool.Get()
	if err != nil {
		DefaultLogger.Log("level", "error", "msg", fmt.Sprintf("Error getting Redis connection: %s", err.Error()))
	}
	defer red.pool.Put(conn)

	b, err := json.Marshal(value)

	if conn.Cmd("SET", red.makeKey(key), b).Err != nil {
		DefaultLogger.Log("level", "error", "msg", fmt.Sprintf("Error setting Redis: %s", err.Error()))
		return err
	}

	return err
}

func (red *RedisJsonStore) Get(key string, value interface{}) (bool, error) {

	PromCacheRequests.With(prom.Labels{"store": "quicktime"}).Inc()

	conn, err := red.pool.Get()
	if err != nil {
		DefaultLogger.Log("level", "error", "msg", fmt.Sprintf("Error getting Redis connection: %s", err.Error()))
	}
	defer red.pool.Put(conn)

	resp := conn.Cmd("GET", red.makeKey(key))
	bytes, err := resp.Bytes()
	if err != nil {
		DefaultLogger.Log("level", "error", "msg", fmt.Sprintf("Error querying Redis: %s", err.Error()))
		// Differentiate different kinds of errors
		return false, nil
	}

	json.Unmarshal(bytes, value)

	return true, nil
}

func CreateRedisJSONStore(redisHost, prefix string) (*RedisJsonStore, error) {
	p, err := pool.New("tcp", redisHost, 10)
	if err != nil {
		return &RedisJsonStore{}, err
	}

	return &RedisJsonStore{
		pool:   p,
		prefix: prefix,
	}, nil
}
