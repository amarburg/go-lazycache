package lazycache

import "fmt"
import "net/http"
import "strings"

import "sync"

type Node struct {
	Path, trimPath string
	Children       map[string]*Node
	leafFunc       func(*Node, []string, http.ResponseWriter, *http.Request) *Node
	Fs             *HttpFS

	updateMutex sync.Mutex
}

type RootNode struct {
	node *Node

	nodeMap map[string]*Node
}


var RootMap map[string]*RootNode

func init() {
	RootMap = make(map[string]*RootNode)
}



func (root RootNode) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	DefaultLogger.Log("msg",fmt.Sprintf("In rootNode::ServeHTTP for %s", root.node.Fs.Uri.String() ) )
	// Sanitive the input URL
	shortPath := strings.TrimPrefix(req.URL.Path, root.node.trimPath)
	elements := stripBlankElementsRight(strings.Split(shortPath, "/"))

	// Launch into the tree
	root.Handle(root.node, elements, w, req)
}

func (root *RootNode) Handle(node *Node, path []string, w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Calling Handle for path %s with (%d): (%s)\n", node.Path, len(path), strings.Join(path, ":"))

	// If I have a leafFunc, I've been assigned a Handler.
	// if node.leafFunc == nil {
	// 	// If not, try to autodetect what my job should be
	// 	node.autodetectLeafFunc()
	// }

	if node.leafFunc != nil {
		if recurse := node.leafFunc(node, path, w, req); recurse != nil {
			root.Handle(recurse, path[1:], w, req)
		}
	} else {
		http.Error(w, fmt.Sprintf("Don't know what to do with path %s", node.Path), 400)
	}

	// else if len(path) > 0 {
	// 	// Still no assignment?  If there are paths left, assume it's a directory and recurse
	//
	// 	fmt.Printf("Don't know what to do with %s but there are paths left, assume it's a directory and move on...\n", path[0])
	//
	// 	if _, ok := node.Children[path[0]]; ok == false {
	// 		newNode := node.MakeNode(path[0] + "/")
	// 		node.Children[path[0]] = newNode
	// 		//fmt.Printf("Registering %s\n", newNode.trimPath)
	// 		//http.Handle(newNode.trimPath, newNode)
	// 	}
	//
	// 	//newNode.leafFunc = HandleDirectory
	// 	//newNode.autodetectLeafFunc()
	// 	root.Handle(node.Children[path[0]], path[1:], w, req)
	//
	// } else {

}

func (node *Node) autodetectLeafFunc() {

	if movExtension.MatchString(node.Path) {
		node.leafFunc = MoovHandler
	} else if mp4Extension.MatchString(node.Path) {
		node.leafFunc = CacheHandler
	} else {
		//fmt.Printf("Could not detect type for %s\n", node.Path)
		node.leafFunc = RedirectHandler

	}
}

func (parent *Node) MakeNode(path string) *Node {
	//fmt.Println("Creating node for", path )

	trimPath := parent.trimPath + path
	fullPath := parent.Path + path
	node := &Node{
		Children: make(map[string]*Node),
		Path:     fullPath,
		trimPath: trimPath,
		Fs:       parent.Fs,
	}

	return node
}

func MakeRootNode(Fs *HttpFS, root string) {
	rootNode := &RootNode{
		node: &Node{
			Path:     "/",
			trimPath: root,
			Children: make(map[string]*Node),
			Fs:       Fs,
		},
	}

	http.Handle(rootNode.node.trimPath, rootNode)
	rootNode.node.leafFunc = HandleDirectory

	DefaultLogger.Log("level","debug","msg", fmt.Sprintf("Handling %s", rootNode.node.trimPath ))

	RootMap[Fs.Uri.String()] = rootNode

// Assign leafFunc because we know it's a directory
	//rootNode.node.autodetectLeafFunc()
}
