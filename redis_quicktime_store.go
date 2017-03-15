package lazycache

//import "net/url"

import (
  "fmt"
  "encoding/json"
)

import "github.com/amarburg/go-lazyfs"
import "github.com/amarburg/go-lazyquicktime"

import prom "github.com/prometheus/client_golang/prometheus"

import "github.com/mediocregopher/radix.v2/pool"

const qtPrefix = "qt."

type RedisQuicktimeStore struct {
  pool      *pool.Pool
}

func (red *RedisQuicktimeStore) Update( key string, fs lazyfs.FileSource ) (*lazyquicktime.LazyQuicktime,error) {
  var err error
metadata,err := lazyquicktime.LoadMovMetadata( fs )
if err != nil {
  DefaultLogger.Log("level", "error", "msg", fmt.Sprintf("Error parsing Quicktime metadata for %s: %s", key, err.Error() ) )
return nil, err
}

  // PromCacheSize.With( prom.Labels{"store":"quicktime"}).Set( float64(len(red.store)))

  conn, err := red.pool.Get()
  if err != nil {
  	DefaultLogger.Log("level", "error", "msg", fmt.Sprintf("Error getting Redis connection: %s", err.Error() ) )
  }
  defer red.pool.Put( conn )

  b, err := json.Marshal(metadata)

  if conn.Cmd("SET", qtPrefix + key, b ).Err != nil {
    DefaultLogger.Log("level", "error", "msg", fmt.Sprintf("Error setting Redis: %s", err.Error() ) )
    return nil, err
  }

  return  metadata, err
}


func (red *RedisQuicktimeStore) Get( key string ) (*lazyquicktime.LazyQuicktime, bool) {

  PromCacheRequests.With( prom.Labels{"store":"quicktime"}).Inc()

  conn, err := red.pool.Get()
  if err != nil {
  	DefaultLogger.Log("level", "error", "msg", fmt.Sprintf("Error getting Redis connection: %s", err.Error() ) )
  }
  defer red.pool.Put( conn )

  resp := conn.Cmd("GET", qtPrefix + key )
  bytes,err := resp.Bytes()
  if err != nil {
    DefaultLogger.Log("level", "error", "msg", fmt.Sprintf("Error querying Redis: %s", err.Error() ) )
    // Differentiate different kinds of errors
    return nil, false
  }

fmt.Println( string(bytes) )

 qt := &lazyquicktime.LazyQuicktime{}
 json.Unmarshal( bytes, qt )

 return qt, true
}

func CreateRedisQuicktimeStore( host string ) (*RedisQuicktimeStore, error) {
  p, err := pool.New("tcp", "localhost:6379", 10)
  if err != nil {
    return &RedisQuicktimeStore{}, err
  }

  return &RedisQuicktimeStore{
    pool: p,
  }, nil
}
