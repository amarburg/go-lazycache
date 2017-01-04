package main

import (
  "fmt"
  "net/http"
  "regexp"
  "encoding/json"
)
import "github.com/amarburg/go-lazyfs"

var OOIRawDataRootURL = "https://rawdata.oceanobservatories.org/"

var fs *lazyfs.HttpFSSource = nil
var trailingSlash = regexp.MustCompile(`/$`)

func index(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(w, "<a href=\"rawdata.oceanobservatories.org/\">rawdata.oceanobservatories.org/</a>\n")
}


func rawData( w http.ResponseWriter, req *http.Request ) {
  listing,err := fs.ReadHttpDir( req.URL.Path )

  if err == nil {
      for _,child := range listing.Children {
        if( trailingSlash.MatchString( child ) ) {
        fmt.Fprintf(w, "<a href=\"%s\">%s</a><br>\n", child, child )
      } else {
        fmt.Fprintf(w, "%s\n", child )
      }
    }

    fmt.Fprintf(w, "\n")
  } else {
    fmt.Fprintf( w, "Error: %s\n", err.Error() )
  }

}


func rawDataJson( w http.ResponseWriter, req *http.Request ) {
  listing,err := fs.ReadHttpDir( req.URL.Path )

  if err == nil {

    type DirListing struct {
      Directories []string
      Files []string
    }

    output := DirListing{}

    for _,child := range listing.Children {
      if( trailingSlash.MatchString( child ) ) {
        output.Directories = append(output.Directories, child )
      } else {
        output.Files = append(output.Files, child )
      }
    }


    b, err := json.MarshalIndent(output,"","  ")
    if err != nil {
      fmt.Fprintln(w, "JSON error:", err)
    }

    w.Write(b)

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

  fmt.Printf("Starting http handler at http://localhost:5000/\n")

  http.HandleFunc("/", rawDataJson)
  //http.HandleFunc("/", index)
  http.ListenAndServe("127.0.0.1:5000", nil)
}
