package lazycache

// import "cloud.google.com/go/storage"
// import "golang.org/x/net/context"
import "io"
// import "fmt"


type ImageStore interface {
	 Has(key string) bool
	 Url(key string) (string, bool)
	 Store(key string, data io.Reader)
	 Retrieve(key string) (io.Reader, error)
}

var DefaultImageStore ImageStore = NullImageStore{}
//
// // Singletons which wrap the
// func Has(key string) bool {
// 	return DefaultImageStore.Has(key)
// }
//
// func Store(key string, data io.Reader) {
// 	DefaultImageStore.Store(key, data)
// }
//
// func Retrieve(key string) (io.Reader, error) {
// 	return DefaultImageStore.Retrieve(key)
// }
//
// func Url(key string) (string, bool) {
// 	return DefaultImageStore.Url(key)
// }
