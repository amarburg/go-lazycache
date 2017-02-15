package lazycache

import (
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
  "golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"io"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
)


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

func (store GoogleImageStore) Statistics() ( interface {} ) {
	return struct{
			Type string
		}{
			Type: "google_cloud_storage",
		}
}

func CreateGoogleStore( bucket string, logger kitlog.Logger  ) (GoogleImageStore){
	logger = kitlog.NewContext(logger).With("module", "GoogleImageStore")

  store := GoogleImageStore{}

  fmt.Printf("Creating Google image store in bucket \"%s\"\n", bucket)

	var err error
	store.ctx = context.Background()

	cred,err := google.FindDefaultCredentials( store.ctx, "https://www.googleapis.com/auth/cloud-platform" )
	if err != nil {
		logger.Log("level","info",
								"tag","auth_error",
								"msg",fmt.Sprintf("Credential error: %s", err.Error()))

	}
	store.client, err = storage.NewClient(store.ctx,
																				option.WithTokenSource(cred.TokenSource) )
	if err != nil {
		panic(fmt.Sprintf("Error opening storage client: %s", err.Error()))
	}

	store.bucket = store.client.Bucket( bucket )
	// if err := store.bucket.Create(store.ctx, ProjectId, nil); err != nil {
	//   panic(fmt.Sprintf("Error creating bucket: %s", err.Error()))
	// }

return store

}
