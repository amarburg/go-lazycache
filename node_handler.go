package lazycache

import "fmt"
import "net/http"
import "strings"

import "regexp"

type Node struct {
  Path, trimPath    string
  Fs *HttpFS
  Children map[string]*Node
  leafFunc func( *Node, []string, http.ResponseWriter, *http.Request )
}

var movExtension = regexp.MustCompile(`\.mov$`)
var mp4Extension = regexp.MustCompile(`\.mp4$`)

// // type RootNode struct {
// //   nodes    *Node
// //   rootPath string
// // }

func stripBlankElementsRight( slice []string ) []string {
  if len(slice) > 0 && len( slice[len(slice)-1] ) == 0  {
      return stripBlankElementsRight( slice[:len(slice)-1] )
  }
  return slice
}

func (node Node) ServeHTTP( w http.ResponseWriter, req *http.Request ) {

  shortPath := strings.TrimPrefix( req.URL.Path, node.trimPath )
  elements  := stripBlankElementsRight( strings.Split( shortPath, "/" ) )

  fmt.Printf("ServeHTTP with %d elements: (%s)\n", len(elements), strings.Join( elements, ":"))

  // Starting root, pass off to Handlers
  node.Handle( elements, w, req )
}


func (node *Node) Handle( path []string, w http.ResponseWriter, req *http.Request ) {
  fmt.Printf("Calling Handler %s with path (%d): (%s)\n", node.Path, len(path), strings.Join(path,":") )

  // If I have a leafFunc, I've been assigned a Handler.
  if node.leafFunc == nil {
    // If not, try to autodetect what my job should be
    node.autodetectLeafFunc()
  }

  if node.leafFunc != nil {
    // args := path
    // if len(path) > 0 {
    //   args = path[1:]
    // }
    node.leafFunc( node, path, w, req )
  } else {
    // Still no assignment?  If there are paths left, assume it's a directory and recurse

    fmt.Printf("Don't know what to do with %s but there are paths left, assume it's a directory and move on...", path[0])
    if len(path) > 0 {
      newNode := node.MakeNode( path[0] + "/" )
      newNode.leafFunc = HandleDirectory
      node.Children[path[0]] = newNode
      newNode.autodetectLeafFunc()
      newNode.Handle( path[1:], w, req )
    }
  }

}

  // if( len( path ) > 0 ) {
  //   fmt.Printf( "Node handling: %s\n", node.Path )
  //
  //   if child,ok := node.children[ path[0] ]; ok   {
  //     if child != nil {
  //       child.Handle( path[1:], w, req )
  //     } else {
  //
  //       // Create a new directory node, populate it, then run it
  //       newNode := node.makeNode( path )
  //       node.children[ path[0] ] = newNode
  //       newNode.populate()
  //       newNode.Handle( path[1:], w, req )
  //
  //     }
  //   } else {
  //     fmt.Println("New path: ")
  //   }
  // } else {
  //   fmt.Println("len(path) == 0")
  //   // printf("Leaf: %s \n", node.Path)
  //
  // }
//}


func (node *Node) autodetectLeafFunc() {

  if movExtension.MatchString( node.Path ) {
    node.leafFunc = HandleMov
  } else if mp4Extension.MatchString( node.Path ) {
    node.leafFunc = HandleDefault
  } else {
    // Try a directory

    listing,err := node.Fs.ReadHttpDir( node.Path )

    if err == nil {
      node.leafFunc = HandleDirectory

      fmt.Printf("Auto detected a directory...\n")
      BootstrapDirectory( node, listing )
      // TODO.  Reformat the output for JSON
      //fmt.Printf("Populating node %s with %d children and %d files\n", node.Path, len(listing.Files), len(listing.Directories))

      // for _,d := range listing.Directories {
      //   node.children[ d ] = nil
      // }

    } else {
      fmt.Printf("Could not detect type for %s\n", node.Path )
    }

  }

}


func MakeRootNode( Fs *HttpFS, root string ) (*Node) {
  node := &Node{Path: "/",
                trimPath: root,
                Children: make( map[string]*Node ),
                Fs: Fs,
              }

  node.autodetectLeafFunc()

  fmt.Println("registering ", node.trimPath )
  http.Handle( node.trimPath, node )

  return node
}

func (parent *Node) MakeNode( path string ) (*Node) {
  fmt.Println("Creating node for", path )

  trimPath := parent.trimPath + path
  fullPath := parent.Path + path
  node := &Node{ Fs: parent.Fs,
                Children: make( map[string]*Node ),
                Path: fullPath,
                trimPath: trimPath }

  // Assign leafFunc
  // switch( node.Fs.PathType( node.trimPath ) ) {
  // case Directory: node.leafFunc = HandleDirectory
  // case File: if movExtension.MatchString( path[0] ) {
  //               node.leafFunc = HandleMov
  //             }
  // }

  fmt.Println("registering ", node.trimPath )
  http.Handle( node.trimPath, node )

  return node
}




// func (root RootNode) ServeHTTP( w http.ResponseWriter, req *http.Request ) {
//   fmt.Printf("Default Handler %s\n", req.URL.Path )
//
//   shortPath := strings.TrimPrefix( req.URL.Path, root.rootPath )
//   elements  := strings.Split( shortPath, "/" )
//   elements  = elements[:len(elements)-1]
//   // Starting root, pass off to Handlers
//
//   root.nodes.Handle( elements, w, req )
// }
