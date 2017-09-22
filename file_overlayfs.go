package lazycache

import (
	// "bytes"
	// "errors"
	"fmt"
	"github.com/amarburg/go-lazyfs"
	// "github.com/valyala/fasthttp"
	// "golang.org/x/net/html"
	// "net/url"
	// "path"
	// "regexp"
	"os"
	"path/filepath"
)

//====

// type FileSystem interface {
// 	PathType(path string) int
// 	ReadDir(p string) (*DirListing, error)
// 	OriginalPath(p string) string
// 	FileSource(p string) (lazyfs.FileSource, error)
// }

type FileOverlayFS struct {
	fs 				FileSystem
	path			string
	Flatten			bool
}


func OpenFileOverlayFS(fs FileSystem, path string) (*FileOverlayFS, error) {

	ofs := &FileOverlayFS{
				fs: fs,
				path: path,
				Flatten: false,
	}

	return ofs, nil
}

func (fs *FileOverlayFS) OriginalPath(p string) string {
	return fs.fs.OriginalPath(p)
}

func (fs *FileOverlayFS) PathType(path string) int {
	return fs.fs.PathType( path )
}

func (fs *FileOverlayFS) FileSource(p string) (lazyfs.FileSource, error) {

	var localFileName string
	if fs.Flatten {
		localFileName = filepath.Base(p)
	} else {
		localFileName = p
	}

	localPath := filepath.Join( fs.path, localFileName )

	Logger.Log("debug", fmt.Sprintf("Checking file overlay for %s", localPath))

	_,err := os.Stat( localPath )

	if err != nil {
		return fs.fs.FileSource(p)
	}

	Logger.Log("msg", fmt.Sprintf("Using %s as local overlay file for %s", localPath, p ) )

	return lazyfs.OpenLocalFile( localPath )
}

// func (fs *HttpFS ) Open( path string ) (*HttpSource, error) {
//   url,_ := url.Parse(fs.url_root + path)
//   src,err := OpenHttpSource( *url )
//
//   return src,err
// }

func (fs *FileOverlayFS) ReadDir(p string) (*DirListing, error) {
	return fs.fs.ReadDir(p)
}
