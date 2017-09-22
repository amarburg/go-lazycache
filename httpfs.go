package lazycache

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/amarburg/go-lazyfs"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/html"
	"net/url"
	"path"
	"regexp"
)

//====

type FileSystem interface {
	PathType(path string) int
	ReadDir(p string) (*DirListing, error)
	OriginalPath(p string) string
	LazyFile(p string) (lazyfs.FileSource, error)
}

type HttpFS struct {
	Uri    url.URL
	client fasthttp.Client
}

const (
	Directory = iota
	File      = iota
)

var trailingSlash = regexp.MustCompile(`/$`)

func OpenHttpFS(uri url.URL) (*HttpFS, error) {
	fs := HttpFS{Uri: uri}
	return &fs, nil
}

func (fs *HttpFS) PathType(path string) int {

	fmt.Println("Finding pathType of ", path)

	// Pure heuristics right now
	if trailingSlash.MatchString(path) {
		return Directory
	}

	return File
}

func (fs *HttpFS) OriginalUri(p string) url.URL {
	uri := fs.Uri
	uri.Path += p
	return uri
}

func (fs *HttpFS) OriginalPath(p string) string {
	url := fs.OriginalUri(p)
	return url.String()
}

func (fs *HttpFS) LazyFile(p string) (lazyfs.FileSource, error) {
	lazy, err := lazyfs.OpenHttpSource(fs.OriginalUri(p))
	return lazy, err
}

// func (fs *HttpFS ) Open( path string ) (*HttpSource, error) {
//   url,_ := url.Parse(fs.url_root + path)
//   src,err := OpenHttpSource( *url )
//
//   return src,err
// }

func (fs *HttpFS) ReadDir(p string) (*DirListing, error) {

	pathUri := fs.Uri
	pathUri.Path += p

	fmt.Printf("Querying directory: %s\n", pathUri.String())

	statusCode, body, err := fs.client.Get([]byte{}, pathUri.String())

	listing := &DirListing{Path: p,
		Files:       make([]string, 0),
		Directories: make([]string, 0),
	}

	//fmt.Println( response, err )

	if err != nil {
		return listing, err
	} else if statusCode != 200 {
		return listing, errors.New(fmt.Sprintf("Got HTTP response %d", statusCode))
	}

	d := html.NewTokenizer(bytes.NewBuffer(body))

	//fmt.Println(d)

	for {
		tt := d.Next()

		if tt == html.ErrorToken {
			break
		}

		// fmt.Println(tt)

		// Big ugly brutal
		switch tt {
		case html.StartTagToken:
			token := d.Token()

			for _, attr := range token.Attr {
				if attr.Key == "href" {
					//fmt.Println(attr.Key)
					val := attr.Val

					tt = d.Next()
					if tt == html.TextToken {
						next := d.Token()
						text := next.Data
						//fmt.Println(text)

						if val == text {

							if trailingSlash.MatchString(text) {
								listing.Directories = append(listing.Directories, path.Clean(text))
							} else {
								listing.Files = append(listing.Files, text)
							}

						}
					}
				}
			}

		}
	}

	return listing, err
}
