package main

import "os"
import "fmt"
import "strings"
import "menteslibres.net/gosexy/yaml"


const (
  IMAGE_STORE_NONE = iota
  IMAGE_STORE_GOOGLE = iota
)

type ImageStoreConfig struct {
  Type   int
  Bucket  string
}

type LazyCacheConfig struct {
  RootUrl      string
  ServerIp      string
  ServerPort    int
  ImageStoreConfig    ImageStoreConfig
}

const ImageStoreConfigName = "image_store"

func LoadLazyCacheConfig( filename string ) (LazyCacheConfig, error) {
  // Default values
  config := LazyCacheConfig{
    ServerPort:   5000,
    ServerIp:     "0.0.0.0",
  }

  _,err := os.Stat( filename )
  if len(filename) == 0 || os.IsNotExist( err ) {
    return config, fmt.Errorf( "Need to specify a valid configuration with the --config option." )
  }

settings, err := yaml.Open(filename)



  if settings == nil { return config, err }

  fmt.Println(settings)

  val := settings.Get("server_ip")
  if val != nil { config.ServerIp = val.(string) }

  val = settings.Get("server_port")
  if val != nil { config.ServerPort = val.(int) }

  //-- Configure root_uri
  val = settings.Get("root_url")
  if val == nil {
    return config, fmt.Errorf("RootURL not specified")
  }
  config.RootUrl = val.(string)


  //-- Configure image store
  _,err = config.ImageStoreConfig.loadConfig( settings )
  if err != nil { return config, err }


  return config, nil
}


func (config *ImageStoreConfig) loadConfig( settings *yaml.Yaml ) (bool,error) {

  // Set defaults
config.Type = IMAGE_STORE_NONE

val := settings.Get(ImageStoreConfigName)
if val == nil {  return true, nil }

val = settings.Get( ImageStoreConfigName, "type" )
if val == nil { return true, nil }

fmt.Printf("ImageStore type: %s", val)
switch( strings.ToLower( val.(string) )) {
case "google":
  config.Type = IMAGE_STORE_GOOGLE
default:
  return false, fmt.Errorf("Don't recognize ImageStore type \"%s\"", val )
}

switch( config.Type ) {
case IMAGE_STORE_GOOGLE:
  bucket := settings.Get( ImageStoreConfigName, "bucket")
  if bucket == nil {
    return false, fmt.Errorf("Google Image store requires a bucket name")
  }
  config.Bucket = bucket.(string)
}

  return true, nil
}

// func parseLazyCacheConfig( hash map[string]string ) LazyCacheConfig {
//   config := LazyCacheConfig{}
//
//   var ok bool
//   config.RootUrl, ok = hash["RootUrl"]
//   if !ok {
//     panic("RootUrl is required in the config file")
//   }
//
//   imageStore,ok := hash["ImageStore"]
//   fmt.Println(imageStore)
//   // if ok {
//   //   switch imageStore["Type"] {
//   //   case "GoogleBucket":
//   //     fmt.Println("Using a Google Bucket for image store")
//   //   default:
//   //     panic(fmt.Sprintf("Don't know how to make an ImageStore of type \"%s\"", imageStore["Type"]))
//   //   }
//   // }
//
//
//
//   return config
// }
