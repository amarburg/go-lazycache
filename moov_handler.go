package lazycache

import "fmt"
import "net/http"

func HandleMov( node *Node, path []string, w http.ResponseWriter, req *http.Request ) {
  fmt.Fprintf( w, "Quicktime handler: %s\n", node.Path )
}


func HandleDefault( node *Node, path []string, w http.ResponseWriter, req *http.Request ) {
  fmt.Fprintf( w, "Default handler: %s\n", node.Path )
}
