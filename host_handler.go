package lazycache

import "fmt"
import "net/http"
import "strings"

type HandlerCommon struct {
//  root string   // For simplicity, root is stored without the trailing slash
  fs *HttpFS
}

type Node struct {
  path, trimPath    string
  common *HandlerCommon
  children map[string]*Node
}

// type RootNode struct {
//   nodes    *Node
//   rootPath string
// }



func (node *Node) handle( path []string, w http.ResponseWriter, req *http.Request ) {
  fmt.Fprintf( w, "Node: %s\n", node.path )

  if( len( path ) > 0 ) {
    child,ok := node.children[ path[0] ]
    if ok  {
      child.handle( path[1:], w, req )
    } else {
      newNode := node.makeNode( path )
      node.children[ path[0] ] = newNode
      newNode.handle( path[1:], w, req )
    }
  }
}

func (handle *Node) makeNode( path []string ) (*Node) {
  fmt.Println("Creating node for", path[0] )

  fullPath := handle.trimPath + path[0] + "/"
  node := Node{ common: handle.common,
                children: make( map[string]*Node ),
                path: path[0] + "/",
                trimPath: fullPath }

  fmt.Println("registering ", fullPath )
  http.Handle( fullPath, node )

  return &node
}

func (node Node) ServeHTTP( w http.ResponseWriter, req *http.Request ) {
  fmt.Printf("Node handler %s for %s\n", node.path, req.URL.Path )

  shortPath := strings.TrimPrefix( req.URL.Path, node.trimPath )
  elements  := strings.Split( shortPath, "/" )
  elements  = elements[:len(elements)-1]

  // Starting root, pass off to handlers
  node.handle( elements, w, req )
}


func MakeRootNode( fs *HttpFS, root string ) (*Node) {
  return &Node{path: "/",
                trimPath: root,
                children: make( map[string]*Node ),
                common: &HandlerCommon{ fs: fs },
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
