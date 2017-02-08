package main

//import "fmt"
import "net/http"

func HandleCache(node *Node, path []string, w http.ResponseWriter, req *http.Request) *Node {
	//fmt.Fprintf( w, "Redirect handler: %s\n", node.Path )
	cacheUrl := node.Fs.Uri
	cacheUrl.Path += node.Path
	http.Redirect(w, req, cacheUrl.String(), 302)

	return nil
}
