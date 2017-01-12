package main

import "fmt"
import "net/http"
import "net/url"

import "github.com/amarburg/go-lazycache"

var OOIRawDataRootURL = "https://rawdata.oceanobservatories.org/"

func main() {

  url,err := url.Parse( OOIRawDataRootURL )
  fs, err := lazycache.OpenHttpFS( *url )

  if err != nil {
    panic( fmt.Sprintf("Error opening HTTP FS Source: %s", err.Error() ) )
  }

  serverAddr := "localhost:5000"


  //http.HandleFunc("*.mov/*", lazycache.MoovHandler )

  root := fmt.Sprintf("%s/%s/", fs.Uri.Host, fs.Uri.Path )
  http.Handle(root, lazycache.MakeTreeHandler( fs ) )
  //http.HandleFunc("/", index)

  fmt.Printf("Starting http handler at http://%s/\n", serverAddr)
  fmt.Printf("Fs at http://%s/%s\n", serverAddr, root )

  http.ListenAndServe(serverAddr, nil)
}
