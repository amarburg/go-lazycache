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
	//  Statistics() ( interface {} )
}

var ImageCache ImageStore = NullImageStore{}

//
// // Singletons which wrap the
// func Has(key string) bool {
// 	return ImageCache.Has(key)
// }
//
// func Store(key string, data io.Reader) {
// 	ImageCache.Store(key, data)
// }
//
// func Retrieve(key string) (io.Reader, error) {
// 	return ImageCache.Retrieve(key)
// }
//
// func Url(key string) (string, bool) {
// 	return ImageCache.Url(key)
// }
