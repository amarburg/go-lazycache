package lazycache

import "net/http"
import "fmt"
import "encoding/json"
import "strings"
import "time"

import "sync"

type DirListing struct {
	Path        string
	Files       []string
	Directories []string
	expires			time.Time
}

type DirMapStore struct {
	Cache map[string](*DirListing)
	Mutex sync.Mutex
}

var DirCache DirMapStore

// a Duration, measured in nanoseconds
const CachedExpiration = 5 * time.Minute

func init() {
	DirCache = DirMapStore{
		Cache: make(map[string](*DirListing)),
	}
}

func (cache *DirMapStore) getDirectory(node *Node) (*DirListing, error) {

	// Initialize or update as necessary
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	cacheKey := node.Fs.OriginalPath(node.Path)
	//Logger.Log("msg", fmt.Sprintf("Checking cache key: %s", cacheKey))

	listing, has := cache.Cache[cacheKey]

	if has {
		if time.Now().After( listing.expires ) {
			has = false
		}
	}

	if !has {
		//Logger.Log("msg", fmt.Sprintf("Need to update dir cache for %s", node.Path))
		var err error
		listing, err = node.Fs.ReadDir(node.Path)

		//Logger.Log("msg", fmt.Sprintf("Listing has %d files and %d directories", len(listing.Files), len(listing.Directories)))

		if err == nil {
			listing.expires = time.Now().Add( CachedExpiration )
			//Logger.Log( "msg", fmt.Sprintf("Created cache entry for %s, current time is %s,  expires at %s", cacheKey, time.Now(), listing.expires.String() ))
			cache.Cache[cacheKey] = listing
		} else {
			Logger.Log("msg", fmt.Sprintf("Error querying remote directory: %s", node.Path))
		}
	}

	// This needs to be inside the mutex because it changes listing
	// How else can I tell if the node tree needs to be updated?
	if len(node.Children) != len(listing.Directories)+len(listing.Files) {
		Logger.Log("msg", fmt.Sprintf("Bootstrapping directory %s (%d != %d+%d)", node.Path,
			len(node.Children), len(listing.Directories), len(listing.Files)))
		node.BootstrapDirectory(*listing)
	}

	return listing, nil
}

func HandleDirectory(node *Node, path []string, w http.ResponseWriter, req *http.Request) *Node {
	//fmt.Printf("HandleDirectory %s with path (%d): (%s)\n", node.Path, len(path), strings.Join(path, ":"))

	listing, err := DirCache.getDirectory(node)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving directory: %s", err.Error()), 500)
		return nil
	}

	//Logger.Log("msg", fmt.Sprintf("Listing has %d files and %d directories\"", len(listing.Files), len(listing.Directories)))

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

		b, err := json.MarshalIndent(listing, "", "  ")
		if err != nil {
			Logger.Log("level", "error", "msg", fmt.Sprintf("JSON error:", err))
		}

		addCacheDefeatHeaders(w)

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
