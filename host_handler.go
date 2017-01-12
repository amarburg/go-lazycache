package lazycache

import "fmt"
import "net/http"
import "strings"

type HandlerCommon struct {
  root string   // For simplicity, root is stored without the trailing slash
  fs *HttpFS
}

type Node struct {
  path, fullPath    string
  common *HandlerCommon
  children map[string]*Node
}

type RootNode struct {
  nodes    *Node
  rootPath string
}



func (node *Node) Handle( path []string, w http.ResponseWriter, req *http.Request ) {
  fmt.Fprintf( w, "Node: %s\n", node.path )

  if( len( path ) > 0 ) {
    child,ok := node.children[ path[0] ]
    if ok  {
      child.Handle( path[1:], w, req )
    } else {
      node.children[ path[0] ] = node.makeNode( path )
      node.children[ path[0] ].Handle( path[1:], w, req )
    }
  }
}

func (handle *Node) makeNode( path []string ) (*Node) {
  fmt.Println("Creating node for", path[0] )

  fullPath := handle.fullPath + path[0] + "/"
  node := Node{ common: handle.common,
                children: make( map[string]*Node ),
                path: path[0] + "/",
                fullPath: fullPath }

  fmt.Println("registering ", fullPath )
  http.Handle( fullPath, node )

  return &node
}

func (node Node) ServeHTTP( w http.ResponseWriter, req *http.Request ) {
  fmt.Printf("Default handler %s\n", req.URL.Path )

  shortPath := strings.TrimPrefix( req.URL.Path, node.path )
  elements  := strings.Split( shortPath, "/" )
  elements  = elements[:len(elements)-1]
  // Starting root, pass off to handlers

  node.Handle( elements, w, req )
}


func MakeRootNode( fs *HttpFS, root string ) (*RootNode) {
  return &RootNode{
              rootPath: root,
              nodes:  &Node{path: "/",
                          fullPath: "/",
                          children: make( map[string]*Node ),
                          common: &HandlerCommon{ fs: fs,
                                                  root: strings.TrimRight(root,"/"), },
                            },
  }
}



func (root RootNode) ServeHTTP( w http.ResponseWriter, req *http.Request ) {
  fmt.Printf("Default handler %s\n", req.URL.Path )

  shortPath := strings.TrimPrefix( req.URL.Path, root.rootPath )
  elements  := strings.Split( shortPath, "/" )
  elements  = elements[:len(elements)-1]
  // Starting root, pass off to handlers

  root.nodes.Handle( elements, w, req )
}
