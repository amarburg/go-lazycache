package image_store

import "cloud.google.com/go/storage"
import "golang.org/x/net/context"
import "io"
import "fmt"


const ImageBucket = "image-store-test"
const ProjectId = "smiling-gasket-155322"



type ImageStore struct {
  ctx       context.Context
  client    *storage.Client
  bucket    *storage.BucketHandle
//  index  ImageStoreMap
}

var DefaultImageStore = ImageStore {}

// Singletons which wrap the
func Has( key string ) bool {
  return DefaultImageStore.Has( key )
}

func Store( key string, data io.Reader ) {
  DefaultImageStore.Store( key, data )
}

func Retrieve( key string )  (io.Reader,error)  {
  return DefaultImageStore.Retrieve( key )
}


func Url( key string ) (string, bool) {
  return DefaultImageStore.Url( key)
}



func (store *ImageStore) Has( key string ) (bool) {
  store.initClient()

  return false
}

func (store *ImageStore) Url( key string ) (string, bool) {
  store.initClient()

  obj := store.bucket.Object( key )
  attr,err := obj.Attrs( store.ctx )

  if err != nil {
    return "", false
  }

  return attr.MediaLink, true
}


func (store *ImageStore) Store( key string, data io.Reader ) {
  store.initClient()

  obj := store.bucket.Object( key )
  w := obj.NewWriter(store.ctx)
  io.Copy( w, data )
  if err := w.Close(); err != nil {
      // TODO: Handle error.
  }
}



func (store *ImageStore) Retrieve( key string ) ( io.Reader,error ) {
  store.initClient()

  obj := store.bucket.Object( key )
  reader,err := obj.NewReader(store.ctx)
  if err != nil {
      // TODO: Handle error.
  }

  return reader,err

  // io.Copy( writer, reader )
  // if err := reader.Close(); err != nil {
  //     // TODO: Handle error.
  // }
}



func (store *ImageStore) initClient() {
  if store.client != nil { return }

  var err error
  store.ctx = context.Background()
  store.client, err = storage.NewClient(store.ctx)
  if err != nil {
    panic(fmt.Sprintf("Error opening storage client: %s", err.Error()))
  }

  store.bucket = store.client.Bucket( ImageBucket )
  // if err := store.bucket.Create(store.ctx, ProjectId, nil); err != nil {
  //   panic(fmt.Sprintf("Error creating bucket: %s", err.Error()))
  // }
}
