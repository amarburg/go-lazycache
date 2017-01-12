package lazycache

import (
  "fmt"
  "net/http"
  "encoding/json"
)

type TreeHandler struct {
  fs *HttpFS
}

func MakeTreeHandler( fs *HttpFS ) (TreeHandler) {
  return TreeHandler{ fs: fs }
}

func (handler *TreeHandler) HandleDirectory( w http.ResponseWriter, req *http.Request ) {
  listing,err := handler.fs.ReadHttpDir( req.URL.Path )

  if err == nil {

    // TODO.  Reformat the output for JSON

    b, err := json.MarshalIndent(listing,"","  ")
    if err != nil {
      fmt.Fprintln(w, "JSON error:", err)
    }

    w.Write(b)

  } else {
    http.Error( w, fmt.Sprintf( "Error: %s", err.Error() ), 500 )
  }
}

func (handler TreeHandler) ServeHTTP( w http.ResponseWriter, req *http.Request ) {

  switch( handler.fs.PathType( req.URL.Path) ) {
  case Directory:  handler.HandleDirectory( w, req )
  default: http.Error( w, fmt.Sprintf("Didn't know what to do with: %s", req.URL.Path), 500 )
  }

}


// func index(w http.ResponseWriter, req *http.Request) {
//   fmt.Fprintf(w, "<a href=\"rawdata.oceanobservatories.org/\">rawdata.oceanobservatories.org/</a>\n")
// }


// func rawData( w http.ResponseWriter, req *http.Request ) {
//   listing,err := Fs.ReadHttpDir( req.URL.Path )
//
//   if err == nil {
//       for _,child := range listing.Children {
//         if( trailingSlash.MatchString( child ) ) {
//         fmt.Fprintf(w, "<a href=\"%s\">%s</a><br>\n", child, child )
//       } else {
//         fmt.Fprintf(w, "%s\n", child )
//       }
//     }
//
//     fmt.Fprintf(w, "\n")
//   } else {
//     fmt.Fprintf( w, "Error: %s\n", err.Error() )
//   }
//
// }
