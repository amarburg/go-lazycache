package lazycache

import (
	"cloud.google.com/go/storage"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"io"
)

type GoogleImageStore struct {
	ctx    context.Context
	client *storage.Client
	bucket *storage.BucketHandle
	//  index  ImageStoreMap

	Stats struct {
		cacheRequests int
		cacheMisses   int
	}
}

func (store GoogleImageStore) Has(key string) bool {
	return false
}

func (store *GoogleImageStore) Url(key string) (string, bool) {

	obj := store.bucket.Object(key)
	attr, err := obj.Attrs(store.ctx)

	store.Stats.cacheRequests++

	if err != nil {
		store.Stats.cacheMisses++
		return "", false
	}

	return attr.MediaLink, true
}

func (store GoogleImageStore) Store(key string, data io.Reader) {

	obj := store.bucket.Object(key)
	w := obj.NewWriter(store.ctx)
	io.Copy(w, data)
	if err := w.Close(); err != nil {
		fmt.Printf("Error storing key %s to bucket: %s\n", key, err.Error())
	}

	_, err := obj.Update(store.ctx, storage.ObjectAttrsToUpdate{
		ContentDisposition: "attachment",
		ACL: []storage.ACLRule{
			storage.ACLRule{ Entity: storage.AllUsers,
											 Role: storage.RoleReader,
			},
		},
	})
	if err != nil {
		fmt.Printf("Error setting attributes on %s: %s\n", key, err.Error())
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

func (store GoogleImageStore) Statistics() interface{} {
	return struct {
		Type          string
		CacheRequests int `json: "cache_requests"`
		CacheMisses   int `json: "cache_misses"`
	}{
		Type:          "google_cloud_storage",
		CacheRequests: store.Stats.cacheRequests,
		CacheMisses:   store.Stats.cacheMisses,
	}
}

func CreateGoogleStore(bucket string) *GoogleImageStore {
	logger := kitlog.With(Logger, "module", "GoogleImageStore")

	store := &GoogleImageStore{}

	fmt.Printf("Creating Google image store in bucket \"%s\"\n", bucket)

	var err error
	store.ctx = context.Background()

	cred, err := google.FindDefaultCredentials(store.ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		logger.Log("level", "info",
			"tag", "auth_error",
			"msg", fmt.Sprintf("Credential error: %s", err.Error()))

	}
	store.client, err = storage.NewClient(store.ctx,
		option.WithTokenSource(cred.TokenSource))
	if err != nil {
		panic(fmt.Sprintf("Error opening storage client: %s", err.Error()))
	}

	store.bucket = store.client.Bucket(bucket)
	attrs := &storage.BucketAttrs{StorageClass: "REGIONAL", Location: "us-central1"}
	if err := store.bucket.Create(store.ctx, "camhd-app-dev", attrs); err != nil {
		logger.Log("msg", fmt.Sprintf("Error during bucket creation: %s", err.Error()))
		//panic(fmt.Sprintf("Error creating bucket: %s", err.Error()))
	}

	return store

}
