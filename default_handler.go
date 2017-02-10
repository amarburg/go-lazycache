package lazycache

import "fmt"
import "net/http"

func HandleDefault(node *Node, path []string, w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Default handler: %s\n", node.Path)
}
