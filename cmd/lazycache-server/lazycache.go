package main

import "fmt"
import "net/http"
import "net/url"
import "strings"

import "github.com/amarburg/go-lazycache"

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
