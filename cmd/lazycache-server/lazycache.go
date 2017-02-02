package main

import "fmt"
import "net/http"
import "net/url"
import "strings"

import "github.com/amarburg/go-lazycache"
import "github.com/amarburg/go-lazycache/image_store"


import "flag"

var configFileFlag = flag.String("config", "", "YAML Configuration file")

var OOIRawDataRootURL = "https://rawdata.oceanobservatories.org/"


func main() {

  flag.Parse()

  config,err := LoadLazyCacheConfig( *configFileFlag )

  if err != nil {
    panic( fmt.Sprintf("Error parsing config: %s", err.Error() ) )
  }

  fmt.Println(config)
  ConfigureImageStore( config )

  url,err := url.Parse( config.RootUrl )
  fs, err := lazycache.OpenHttpFS( *url )

  if err != nil {
    panic( fmt.Sprintf("Error opening HTTP FS Source: %s", err.Error() ) )
  }

  serverAddr := fmt.Sprintf("%s:%d", config.ServerIp, config.ServerPort)

  // Reverse hostname
  splitHN := lazycache.MungeHostname( fs.Uri.Host )

  http.HandleFunc("/", lazycache.Index )

  root := fmt.Sprintf("/%s%s", strings.Join(splitHN,"/"), fs.Uri.Path )
  lazycache.MakeRootNode( fs, root )

  fmt.Printf("Starting http handler at http://%s/\n", serverAddr)
  fmt.Printf("Fs at http://%s%s\n", serverAddr, root )

  http.ListenAndServe(serverAddr, nil)
}

func ConfigureImageStore( config LazyCacheConfig) {
  switch( config.ImageStoreConfig.Type ) {
  case IMAGE_STORE_NONE:
    image_store.DefaultImageStore = image_store.NullImageStore{}
    case IMAGE_STORE_GOOGLE:
      fmt.Printf("Creating Google image store in bucket \"%s\"\n", config.ImageStoreConfig.Bucket)
        image_store.DefaultImageStore = image_store.CreateGoogleStore( config.ImageStoreConfig.Bucket )
  }
}
