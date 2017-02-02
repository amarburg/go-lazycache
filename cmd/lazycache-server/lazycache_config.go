package main

import "os"
import "fmt"
import "menteslibres.net/gosexy/yaml"


const (
  IMAGE_STORE_NONE = iota
  IMAGE_STORE_GOOGLE = iota
)

type ImageStoreConfig struct {
  Type     int
}

type LazyCacheConfig struct {
  RootUrl      string
  ServerIp      string
  ServerPort    int
  ImageStore   ImageStoreConfig
}

func LoadLazyCacheConfig( filename string ) (LazyCacheConfig, error) {
  _,err := os.Stat( filename )
  if len(filename) == 0 || os.IsNotExist( err ) {
    panic( "Need to specify a valid configuration with the --config option." )
  }

settings, err := yaml.Open(filename)

  // Default values
  config := LazyCacheConfig{
    ServerPort:   5000,
    ServerIp:     "0.0.0.0",
  }


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
  config.ImageStore = ImageStoreConfig { Type: IMAGE_STORE_NONE }
  val = settings.Get("image_store")
  if val == nil {
    return config, fmt.Errorf("RootURL not specified")
  }
  _,err = config.ImageStore.loadConfig( val )
  if err != nil { return config, err }


  return config, nil
}


func (config *ImageStoreConfig) loadConfig( foo interface{} ) (bool,error) {

fmt.Println(foo)

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
