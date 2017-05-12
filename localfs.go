package lazycache

import (
	"fmt"
	"github.com/amarburg/go-lazyfs"
	"os"
	"path"
)

//====

type LocalFS struct {
	Path string
}

func OpenLocalFS(path string) (*LocalFS, error) {
	fs := &LocalFS{Path: path} //Uri: uri}
	return fs, nil
}

func (fs *LocalFS) PathType(path string) int {

	fmt.Println("Finding pathType of ", path)

	// Pure heuristics right now
	if trailingSlash.MatchString(path) {
		return Directory
	}

	return File
}

func (fs *LocalFS) OriginalPath(p string) string {
	return path.Join(fs.Path, p)
}

func (fs *LocalFS) LazyFile(p string) (lazyfs.FileSource, error) {
	return lazyfs.OpenLocalFileSource(fs.Path, p)
}

func (fs *LocalFS) ReadDir(p string) (*DirListing, error) {

	listing := &DirListing{Path: p,
		Files:       make([]string, 0),
		Directories: make([]string, 0),
	}

	fullPath := path.Join(fs.Path, p)

	fmt.Printf("Querying directory: %s\n", fullPath)

	dir, err := os.Open(fullPath)
	if err != nil {
		return listing, err
	}

	files, err := dir.Readdir(0)

	if err != nil {
		return listing, err
	}

	for _, f := range files {
		if f.IsDir() {
			listing.Directories = append(listing.Directories, f.Name())
		} else if f.Mode()&(os.ModeDir | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice) == 0 {
			// os.ModeType includes all of the "special types" e.g. pipes, directories, etc., but we allow symlinks...
			listing.Files = append(listing.Files, f.Name())
		}
	}

	return listing, err
}
