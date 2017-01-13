package lazycache

import "net/http"
import "fmt"
import "encoding/json"
import "strings"

func HandleDirectory( node *Node, path []string, w http.ResponseWriter, req *http.Request ) {
  fmt.Printf("HandleDirectory %s with path (%d): (%s)\n", node.Path, len(path), strings.Join(path,":") )

  // In this case, any residual path is further directories...
  if len(path) > 0 {

    if len( node.Children ) == 0  {
      listing,err := node.Fs.ReadHttpDir( node.Path )
      if err == nil {
        BootstrapDirectory( node, listing )
      }
    }

    fmt.Printf("Residual path %d left, recursing to %s\n", len(path), path[0] )
   child,ok := node.Children[ path[0] ]
   //fmt.Println( child, ok )
   if ok && child != nil {
     child.Handle( path[1:], w, req )
   }
  } else {
    // Only dump JSON if you're the leaf node

    listing,err := node.Fs.ReadHttpDir( node.Path )

    if err == nil {

      // Doesn't update ... yet
      // Need to be able to unregister from ServeMux, among other things
      // if len(listing.Directories) + len(listing.Files) != len(node.Children) {
      //   // Updated
      //   fmt.Printf("Updating directory for %s\n", node.Path )
      //   BootstrapDirectory( node, listing )
      // }


      // TODO.  Reformat the output for JSON
      // Technically, I should generate this baed on internal structure, not listing

      b, err := json.MarshalIndent(listing,"","  ")
      if err != nil {
        fmt.Fprintln(w, "JSON error:", err)
      }

      w.Write(b)

    } else {
      http.Error( w, fmt.Sprintf( "Error: %s", err.Error() ), 500 )
    }
  }
}

func BootstrapDirectory( node *Node, listing DirListing ) {
  fmt.Printf("Bootstrapping directory %s\n", node.Path)
  node.Children = make( map[string]*Node )

  for _,d := range listing.Directories {
    // Trim off trailing slash
    dirName := strings.TrimRight( d, "/" )
    newNode := node.MakeNode( dirName + "/" )
    newNode.leafFunc = HandleDirectory
    fmt.Printf("Adding directory %s to %s\n", dirName, node.Path )
    node.Children[dirName] = newNode
  }

  for _,f := range listing.Files {
    newNode := node.MakeNode( f )
    node.Children[f] = newNode
    fmt.Printf("Adding file %s to %s\n", f, node.Path )
  }
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
