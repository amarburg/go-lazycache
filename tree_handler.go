package lazycache

import "net/http"
import "fmt"

func HandleDirectory( node *Node, w http.ResponseWriter, req *http.Request ) {
  fmt.Fprintf( w, "DirectoryHandler: %s\n", node.Path )
}

// import (
//   "fmt"
//   "net/http"
//   "encoding/json"
//   "strings"
// )
//
// type DirHandler struct {
//   common *HandlerCommon
// }
//
// func MakeTreeHandler( fs *HttpFS, root string ) (TreeHandler) {
//   return TreeHandler{ fs: fs, root: strings.TrimRight(root,"/") }
// }
//
// func (handler TreeHandler) ServeHTTP( w http.ResponseWriter, req *http.Request ) {
//
//   fmt.Println("TreeHandler Handling ", req.URL.String() )
//
//   // Strip path off the Request
//   path := strings.TrimPrefix( req.URL.Path, handler.root )
//
//   switch( handler.fs.PathType( path ) ) {
//   case Directory:  handler.HandleDirectory( w, path )
//   default: http.Error( w, fmt.Sprintf("Didn't know what to do with: %s", req.URL.Path), 500 )
//   }
//
// }
//
//
// func (handler *TreeHandler) HandleDirectory( w http.ResponseWriter, path string ) {
//
//   fmt.Println("   Treating", path, "as a directory")
//
//   listing,err := handler.fs.ReadHttpDir( path )
//
//   if err == nil {
//
//     // TODO.  Reformat the output for JSON
//
//     b, err := json.MarshalIndent(listing,"","  ")
//     if err != nil {
//       fmt.Fprintln(w, "JSON error:", err)
//     }
//
//     w.Write(b)
//
//   } else {
//     http.Error( w, fmt.Sprintf( "Error: %s", err.Error() ), 500 )
//   }
// }
//
//
//
//
//
// // func rawData( w http.ResponseWriter, req *http.Request ) {
// //   listing,err := Fs.ReadHttpDir( req.URL.Path )
// //
// //   if err == nil {
// //       for _,child := range listing.Children {
// //         if( trailingSlash.MatchString( child ) ) {
// //         fmt.Fprintf(w, "<a href=\"%s\">%s</a><br>\n", child, child )
// //       } else {
// //         fmt.Fprintf(w, "%s\n", child )
// //       }
// //     }
// //
// //     fmt.Fprintf(w, "\n")
// //   } else {
// //     fmt.Fprintf( w, "Error: %s\n", err.Error() )
// //   }
// //
// // }
