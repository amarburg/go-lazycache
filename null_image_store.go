package lazycache

// import "cloud.google.com/go/storage"
// import "golang.org/x/net/context"
import "io"
import "errors"
// import "fmt"


type NullImageStore struct {

}

func (store NullImageStore) Has(key string) bool {
	return false
}

func (store NullImageStore) Url(key string) (string, bool) {
  return "", false
}

func (store NullImageStore) Store(key string, data io.Reader) {
  //
}

func (store NullImageStore) Retrieve(key string) (io.Reader, error) {
	return nil, errors.New("Cannot retrieve from NullImageStore")
}

func (store NullImageStore) Statistics() ( interface {} ) {
	return struct{}{}
}
