package lazycache

import "fmt"
import "net/http"
import "strings"

import "regexp"


type HandlerCommon struct {
//  root string   // For simplicity, root is stored without the trailing slash
}

type Node struct {
  Path, trimPath    string
  Fs *HttpFS
  children map[string]*Node
  leafFunc func( *Node, http.ResponseWriter, *http.Request )
}

var movExtension = regexp.MustCompile(`\.mov$`)


// type RootNode struct {
//   nodes    *Node
//   rootPath string
// }



func (node *Node) handle( path []string, w http.ResponseWriter, req *http.Request ) {

  if( len( path ) > 0 ) {
  fmt.Printf( "Node: %s\n", node.Path )

    child,ok := node.children[ path[0] ]
    if ok  {
      child.handle( path[1:], w, req )
    } else {
      newNode := node.makeNode( path )
      node.children[ path[0] ] = newNode
      newNode.handle( path[1:], w, req )
    }
  } else {
    printf("Leaf: %s \n", node.Path)
    if node.leafFunc != nil {
      node.leafFunc( node, w, req )
    }
  }
}

func (handle *Node) makeNode( path []string ) (*Node) {
  fmt.Println("Creating node for", path[0] )

  trimPath := handle.trimPath + path[0] + "/"
  fullPath := handle.Path + path[0] + "/"
  node := Node{ Fs: handle.Fs,
                children: make( map[string]*Node ),
                Path: fullPath,
                trimPath: trimPath }

  // Assign leafFunc
  switch( node.Fs.PathType( node.trimPath ) ) {
  case Directory: node.leafFunc = HandleDirectory
  case File: if movExtension.MatchString( path[0] ) {
                node.leafFunc = HandleMov
              }
  }

  fmt.Println("registering ", trimPath )
  http.Handle( trimPath, node )

  return &node
}

func (node Node) ServeHTTP( w http.ResponseWriter, req *http.Request ) {
  fmt.Printf("Node handler %s for %s\n", node.Path, req.URL.Path )

  shortPath := strings.TrimPrefix( req.URL.Path, node.trimPath )
  elements  := strings.Split( shortPath, "/" )
  elements  = elements[:len(elements)-1]

  // Starting root, pass off to handlers
  node.handle( elements, w, req )
}


func MakeRootNode( Fs *HttpFS, root string ) (*Node) {
  return &Node{Path: "/",
                trimPath: root,
                children: make( map[string]*Node ),
                Fs: Fs,
              }
}



// func (root RootNode) ServeHTTP( w http.ResponseWriter, req *http.Request ) {
//   fmt.Printf("Default handler %s\n", req.URL.Path )
//
//   shortPath := strings.TrimPrefix( req.URL.Path, root.rootPath )
//   elements  := strings.Split( shortPath, "/" )
//   elements  = elements[:len(elements)-1]
//   // Starting root, pass off to handlers
//
//   root.nodes.Handle( elements, w, req )
// }
