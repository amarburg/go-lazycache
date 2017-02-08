package main

import "net/http"
import "fmt"
import "encoding/json"
import "strings"

func HandleDirectory(node *Node, path []string, w http.ResponseWriter, req *http.Request) {
	fmt.Printf("HandleDirectory %s with path (%d): (%s)\n", node.Path, len(path), strings.Join(path, ":"))

	// If there's residual path, they must be children (not a verb)
	if len(path) > 0 {

		// If I haven't been initialized, initialize myself as a directory
		if len(node.Children) == 0 {
			if listing, err := node.Fs.ReadHttpDir(node.Path); err == nil {
				BootstrapDirectory(node, listing)
			}
		}

		fmt.Printf("%d elements of residual path left, recursing to %s\n", len(path), path[0])

		if child, ok := node.Children[path[0]]; ok && child != nil {
			child.Handle(path[1:], w, req)
		}
	} else {
		// Only dump JSON if you're the leaf node

		// TODO: Should really dump cached values, not reread from the source
		listing, err := node.Fs.ReadHttpDir(node.Path)

		if err == nil {

			// Doesn't update ... yet
			// Need to be able to unregister from ServeMux, among other things
			// if len(listing.Directories) + len(listing.Files) != len(node.Children) {
			//   // Updated
			//   fmt.Printf("Updating directory for %s\n", node.Path )
			//   BootstrapDirectory( node, listing )
			// }

			// TODO.  Reformat the output for JSON
			// Technically, I should generate this baed on internal structure, not listing

			b, err := json.MarshalIndent(listing, "", "  ")
			if err != nil {
				fmt.Fprintln(w, "JSON error:", err)
			}

			w.Write(b)

		} else {
			http.Error(w, fmt.Sprintf("Error: %s", err.Error()), 500)
		}
	}
}

func BootstrapDirectory(node *Node, listing DirListing) {
	fmt.Printf("Bootstrapping directory %s\n", node.Path)

	// Clear any existing children
	node.Children = make(map[string]*Node)

	for _, d := range listing.Directories {
		// Trim off trailing slash
		dirName := strings.TrimRight(d, "/")
		newNode := node.MakeNode(dirName + "/")
		newNode.leafFunc = HandleDirectory // Assign leafFunc because we know it's a directory
		node.Children[dirName] = newNode
	}

	for _, f := range listing.Files {
		newNode := node.MakeNode(f)
		node.Children[f] = newNode
		fmt.Printf("Adding file %s to %s\n", f, node.Path)
	}
}
