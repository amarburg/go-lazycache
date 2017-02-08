package main

import "net/http"

// Doesn't actually cache yet, just redirects
func CacheHandler(node *Node, path []string, w http.ResponseWriter, req *http.Request) *Node {
	//fmt.Fprintf( w, "Redirect handler: %s\n", node.Path )
	cacheUrl := node.Fs.Uri
	cacheUrl.Path += node.Path
	http.Redirect(w, req, cacheUrl.String(), 302)

	return nil
}

func RedirectHandler(node *Node, path []string, w http.ResponseWriter, req *http.Request) *Node {
	//fmt.Fprintf( w, "Redirect handler: %s\n", node.Path )
	cacheUrl := node.Fs.Uri
	cacheUrl.Path += node.Path
	http.Redirect(w, req, cacheUrl.String(), 302)

	return nil
}
