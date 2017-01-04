package main

import (
  "fmt"
  "net/http"
)
import "github.com/amarburg/go-lazyfs"

var OOIRawDataRootURL = "https://rawdata.oceanobservatories.org/"

var fs *lazyfs.HttpFSSource = nil

func index(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(w, "<a href=\"rawdata.oceanobservatories.org/\">rawdata.oceanobservatories.org/</a>\n")
}


func rawData( w http.ResponseWriter, req *http.Request ) {
  listing,err := fs.ReadHttpDir( req.URL.Path )

  if err == nil {
    for _,child := range listing.Children {
      fmt.Fprintf(w, "<a href=\"%s\">%s</a><br>\n", child, child )
    }

    fmt.Fprintf(w, "\n")
  } else {
    fmt.Fprintf( w, "Error: %s\n", err.Error() )
  }

  }


  func main() {

    var err error
    fs, err = lazyfs.OpenHttpFSSource( OOIRawDataRootURL )

    if err != nil {
      panic( fmt.Sprintf("Error opening HTTP FS Source: %s", err.Error() ) )
    }

    fmt.Printf("Starting http handler at http://localhost:5000/")

    http.HandleFunc("/", rawData)
    //http.HandleFunc("/", index)
    http.ListenAndServe("127.0.0.1:5000", nil)
  }
