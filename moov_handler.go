package lazycache

import "fmt"
import "net/http"
import "strings"

func HandleMov( node *Node, path []string, w http.ResponseWriter, req *http.Request ) {
  fmt.Fprintf( w, "Quicktime handler: %s with residual path (%d): (%s)\n", node.Path, len(path), strings.Join(path,":") )
}


func HandleDefault( node *Node, path []string, w http.ResponseWriter, req *http.Request ) {
  fmt.Fprintf( w, "Default handler: %s\n", node.Path )
}
