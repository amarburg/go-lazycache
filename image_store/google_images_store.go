package image_store

import "cloud.google.com/go/storage"
import "golang.org/x/net/context"
import "io"
import "fmt"


type GoogleImageStore struct {
	ctx    context.Context
	client *storage.Client
	bucket *storage.BucketHandle
	//  index  ImageStoreMap
}

func (store GoogleImageStore) Has(key string) bool {

	return false
}

func (store GoogleImageStore) Url(key string) (string, bool) {

	obj := store.bucket.Object(key)
	attr, err := obj.Attrs(store.ctx)

	if err != nil {
		return "", false
	}

	return attr.MediaLink, true
}

func (store GoogleImageStore) Store(key string, data io.Reader) {

	obj := store.bucket.Object(key)
	w := obj.NewWriter(store.ctx)
	io.Copy(w, data)
	if err := w.Close(); err != nil {
		fmt.Printf("Error storing key %s to bucket: %s\n", key, err.Error() )
	}

	_,err := obj.Update(store.ctx, storage.ObjectAttrsToUpdate{
		ContentDisposition: "attachment",
    ACL:   []storage.ACLRule{
              { storage.AllUsers, storage.RoleReader },
            },
	} )
  if err != nil {
    fmt.Printf("Error setting attributes on %s: %s\n", key, err.Error() )
  }
}

func (store GoogleImageStore) Retrieve(key string) (io.Reader, error) {
	//store.initClient()

	obj := store.bucket.Object(key)
	reader, err := obj.NewReader(store.ctx)
	if err != nil {
		// TODO: Handle error.
	}

	return reader, err

	// io.Copy( writer, reader )
	// if err := reader.Close(); err != nil {
	//     // TODO: Handle error.
	// }
}

func CreateGoogleStore( bucket string ) (GoogleImageStore){

  store := GoogleImageStore{}

  fmt.Printf("Creating Google image store in bucket \"%s\"", bucket)

	var err error
	store.ctx = context.Background()
	store.client, err = storage.NewClient(store.ctx)
	if err != nil {
		panic(fmt.Sprintf("Error opening storage client: %s", err.Error()))
	}

	store.bucket = store.client.Bucket( bucket )
	// if err := store.bucket.Create(store.ctx, ProjectId, nil); err != nil {
	//   panic(fmt.Sprintf("Error creating bucket: %s", err.Error()))
	// }

return store

}
