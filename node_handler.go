package main

import "fmt"
import "net/http"
import "strings"

import "regexp"
import "sync"

type Node struct {
	Path, trimPath string
	Fs             *HttpFS
	ChildrenMutex  sync.Mutex
	Children       map[string]*Node
	leafFunc       func(*Node, []string, http.ResponseWriter, *http.Request)
}

var movExtension = regexp.MustCompile(`\.mov$`)
var mp4Extension = regexp.MustCompile(`\.mp4$`)

func (node Node) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// Sanitive the input URL
	shortPath := strings.TrimPrefix(req.URL.Path, node.trimPath)
	elements := stripBlankElementsRight(strings.Split(shortPath, "/"))

	//fmt.Printf("ServeHTTP with %d elements: (%s)\n", len(elements), strings.Join(elements, ":"))

	// Starting root, pass off to Handlers
	node.Handle(elements, w, req)
}

func (node *Node) Handle(path []string, w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Calling Handle for path %s with (%d): (%s)\n", node.Path, len(path), strings.Join(path, ":"))

	// If I have a leafFunc, I've been assigned a Handler.
	if node.leafFunc == nil {
		// If not, try to autodetect what my job should be
		node.autodetectLeafFunc()
	}

	if node.leafFunc != nil {
		node.leafFunc(node, path, w, req)
	} else if len(path) > 0 {
		// Still no assignment?  If there are paths left, assume it's a directory and recurse

		fmt.Printf("Don't know what to do with %s but there are paths left, assume it's a directory and move on...\n", path[0])

		if len(path) > 0 {
			node.ChildrenMutex.Lock()
			if _, ok := node.Children[path[0]]; ok == false {
				newNode := node.MakeNode(path[0] + "/")
				node.Children[path[0]] = newNode
				fmt.Printf("Registering %s\n", newNode.trimPath)
				http.Handle(newNode.trimPath, newNode)
			}
			node.ChildrenMutex.Unlock()

			//newNode.leafFunc = HandleDirectory
			//newNode.autodetectLeafFunc()
			node.Children[path[0]].Handle(path[1:], w, req)

		} else {
			http.Error(w, fmt.Sprintf("Don't know what to do with path %s", node.Path), 400)
		}
	}

}

func (node *Node) autodetectLeafFunc() {

	if movExtension.MatchString(node.Path) {
		node.leafFunc = MoovHandler
	} else if mp4Extension.MatchString(node.Path) {
		node.leafFunc = HandleCache
	} else {
		// Try to parse it as a directory

		listing, err := node.Fs.ReadHttpDir(node.Path)

		if err == nil {
			node.leafFunc = HandleDirectory

			fmt.Printf("Auto detected a directory...\n")
			BootstrapDirectory(node, listing)

			// TODO.  Reformat the output for JSON
			//fmt.Printf("Populating node %s with %d children and %d files\n", node.Path, len(listing.Files), len(listing.Directories))

			// for _,d := range listing.Directories {
			//   node.children[ d ] = nil
			// }

		} else {
			fmt.Printf("Could not detect type for %s\n", node.Path)
		}

	}

}

func (parent *Node) MakeNode(path string) *Node {
	//fmt.Println("Creating node for", path )

	trimPath := parent.trimPath + path
	fullPath := parent.Path + path
	node := &Node{Fs: parent.Fs,
		Children: make(map[string]*Node),
		Path:     fullPath,
		trimPath: trimPath}

	// By default, don't eager load the children of a new node...

	fmt.Println("registering node at ", node.trimPath)

	return node
}

func MakeRootNode(Fs *HttpFS, root string) *Node {
	node := &Node{Path: "/",
		trimPath: root,
		Children: make(map[string]*Node),
		Fs:       Fs,
	}

	node.autodetectLeafFunc()

	fmt.Println("registering root node at ", node.trimPath)
	http.Handle(node.trimPath, node)

	return node
}
