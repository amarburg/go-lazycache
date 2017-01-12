package lazycache

import "fmt"
import "net/http"
import "net/url"
import "errors"
import "golang.org/x/net/html"
import "regexp"

type HttpFS struct {
  Uri url.URL
}

const (
  Directory = iota
  File      = iota
)

var trailingSlash = regexp.MustCompile(`/$`)


func OpenHttpFS( uri url.URL ) (*HttpFS, error) {
  fs := HttpFS{ Uri: uri }
  return &fs, nil
}

func (fs *HttpFS ) PathType( path string ) (int) {
  // Pure heuristics right now
  if( trailingSlash.MatchString( path ) ) {
    return Directory
  }

  return File
}


// func (fs *HttpFS ) Open( path string ) (*HttpSource, error) {
//   url,_ := url.Parse(fs.url_root + path)
//   src,err := OpenHttpSource( *url )
//
//   return src,err
// }

type DirListing struct {
  Path string
  Files []string
  Directories []string
}

func (fs *HttpFS ) ReadHttpDir( path string ) (DirListing,error){
  client := http.Client{}

  pathUri := fs.Uri
  pathUri.Path += path

fmt.Printf( "Querying: %s\n", pathUri.String() )

  response, err := client.Get( pathUri.String() )

  listing := DirListing{ Path: path }

  //fmt.Println( response, err )

  if( response.StatusCode != 200 ) {
    return listing, errors.New(fmt.Sprintf("Got HTTP response %d", response.StatusCode))
  }

  defer response.Body.Close()
  d := html.NewTokenizer( response.Body )

  //fmt.Println(d)

  for {
  	tt := d.Next()

    if tt == html.ErrorToken { break; }

    // fmt.Println(tt)

    // Big ugly brutal
    switch tt {
  	case html.StartTagToken:
      token := d.Token()

      for _,attr := range token.Attr {
        if attr.Key == "href" {
          //fmt.Println(attr.Key)
          val := attr.Val

          tt = d.Next()
          if( tt == html.TextToken ) {
            next := d.Token()
            text := next.Data
            //fmt.Println(text)

            if val == text  {

              if( trailingSlash.MatchString( text ) ) {
                listing.Directories = append(listing.Directories, text )
              } else {
                listing.Files = append(listing.Files, text )
              }

            }
          }
        }
      }

    }
  }

  return listing, err
}
