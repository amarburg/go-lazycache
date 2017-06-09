package lazycache

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Node struct {
	Path, trimPath string
	Children       map[string]*Node
	leafFunc       func(*Node, []string, http.ResponseWriter, *http.Request) *Node
	Fs             FileSystem
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

	reqStart := time.Now()

	//DefaultLogger.Log("msg", fmt.Sprintf("In rootNode::ServeHTTP for %#v", root))

	// Sanitive the input URL
	shortPath := strings.TrimPrefix(req.URL.Path, root.node.trimPath)
	elements := stripBlankElementsRight(strings.Split(shortPath, "/"))

	// Launch into the tree
	root.Handle(root.node, elements, w, req)

	timeTrack(reqStart, "Full HTTP request")
}

func (root *RootNode) Handle(node *Node, path []string, w http.ResponseWriter, req *http.Request) {
	//fmt.Printf("Calling Handle for path %s with (%d): (%s)\n", node.Path, len(path), strings.Join(path, ":"))

	if node.leafFunc != nil {
		if recurse := node.leafFunc(node, path, w, req); recurse != nil {
			root.Handle(recurse, path[1:], w, req)
		}
	} else {
		http.Error(w, fmt.Sprintf("Don't know what to do with path %s", node.Path), 400)
	}

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

func MakeRootNode(fs FileSystem, root string) {
	rootNode := &RootNode{
		node: &Node{
			Path:     "/",
			trimPath: root,
			Children: make(map[string]*Node),
			Fs:       fs,
		},
	}

	//DefaultLogger.Log("level", "debug",
	//								  "msg", fmt.Sprintf("Adding HTTP %#v handler for %s", rootNode, rootNode.node.trimPath))
	http.Handle(rootNode.node.trimPath, rootNode)
	rootNode.node.leafFunc = HandleDirectory

	DefaultLogger.Log("level", "debug", "msg", fmt.Sprintf("Adding root node %s", rootNode.node.trimPath))

	RootMap[fs.OriginalPath("")] = rootNode

	// Assign leafFunc because we know it's a directory
	//rootNode.node.autodetectLeafFunc()
}
