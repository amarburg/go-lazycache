package lazycache

import "fmt"
import "net/http"
//import "strings"
import "encoding/json"

import "github.com/amarburg/go-lazyfs"
import "github.com/amarburg/go-lazycache/quicktime_store"

func MoovHandler( node *Node, path []string, w http.ResponseWriter, req *http.Request ) {
//  fmt.Fprintf( w, "Quicktime handler: %s with residual path (%d): (%s)\n", node.Path, len(path), strings.Join(path,":") )

  lqt,have := quicktime_store.HaveEntry( node.trimPath )

  if !have {
    uri := node.Fs.Uri
    uri.Path += node.Path
    fmt.Println(uri.String())
    fs,err := lazyfs.OpenHttpSource( uri )
    if err != nil {
      http.Error( w, "Something's went boom opening the HTTP Soruce!", 500 )
      return
    }

    lqt,err = quicktime_store.AddEntry( node.trimPath, fs )
    if err != nil {
      http.Error( w, fmt.Sprintf("Something's went boom parsing the quicktime file: %s", err.Error() ), 500 )
      return
    }
  }

  b, err := json.MarshalIndent(lqt ,"","  ")
  if err != nil {
    fmt.Fprintln(w, "JSON error:", err)
  }

  w.Write(b)

}
