package lazycache

import "fmt"
import "net/http"

func HandleMov( node *Node, w http.ResponseWriter, req *http.Request ) {
  fmt.Fprintf( w, "Quicktime handler: %s\n", node.Path )
}
