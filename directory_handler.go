package lazycache

import "net/http"
import "fmt"
import "encoding/json"
import "strings"

type DirListing struct {
	Path        string
	Files       []string
	Directories []string
}

var DirKeyStore JSONStore

func init() {
	DirKeyStore = CreateMapJSONStore()
}

func HandleDirectory(node *Node, path []string, w http.ResponseWriter, req *http.Request) *Node {
	//fmt.Printf("HandleDirectory %s with path (%d): (%s)\n", node.Path, len(path), strings.Join(path, ":"))

	// Initialize or update as necessary
	DirKeyStore.Lock()
	fmt.Println(DirKeyStore)

	// TODO:  Handler error condition
	listing := &DirListing{}
	ok, _ := DirKeyStore.Get(node.Path, listing)

	if !ok {
		DefaultLogger.Log("msg", fmt.Sprintf("Need to update dir cache for %s", node.Path))
		listing, err := node.Fs.ReadHttpDir(node.Path)
		if err == nil {
			DirKeyStore.Update(node.Path, listing)
			node.BootstrapDirectory(listing)
		} else {
			DefaultLogger.Log(fmt.Sprintf("msg", "Errors querying remote directory: %s", node.Path))
		}
		fmt.Printf("new listing of %s: %v\n", node.Path, listing)
	}
	DirKeyStore.Unlock()

	fmt.Printf("post listing of %s: %v\n", node.Path, listing)

	// If there's residual path, they must be children (not a verb)
	if len(path) > 0 {

		//fmt.Printf("%d elements of residual path left, recursing to %s\n", len(path), path[0])

		if child, ok := node.Children[path[0]]; ok && child != nil {
			return child
		} else {
			http.Error(w, fmt.Sprintf("Can't find %s within %s", path[0], node.trimPath), 404)
		}

	} else {
		// You're the leaf node

		// TODO: Should really dump cached values, not reread from the source
		//listing, err := node.Fs.ReadHttpDir(node.Path)

		// Doesn't update ... yet
		// Need to be able to unregister from ServeMux, among other things
		// if len(listing.Directories) + len(listing.Files) != len(node.Children) {
		//   // Updated
		//   fmt.Printf("Updating directory for %s\n", node.Path )
		//   BootstrapDirectory( node, listing )
		// }

		// TODO.  Reformat the output for JSON
		// Technically, I should generate this based on internal structure, not listing

		b, err := json.MarshalIndent(listing, "", "  ")
		if err != nil {
			DefaultLogger.Log("level", "error", "msg", fmt.Sprintf("JSON error:", err))
		}

		w.Write(b)

	}

	return nil
}

func (node *Node) BootstrapDirectory(listing DirListing) {
	//fmt.Printf("Bootstrapping directory %s\n", node.Path)

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

		newNode.autodetectLeafFunc()

		//fmt.Printf("Adding file %s to %s\n", f, node.Path)
	}
}
