package lazycache

import (
  "testing"
  //"context"
  // "github.com/docker/docker/api/types"
  // "github.com/docker/docker/client"
  testfiles "github.com/amarburg/go-lazyfs-testfiles"
  "github.com/amarburg/go-lazyfs"
  "fmt"
  "reflect"
)

// func StartRedis() {
//   cli, err := client.NewEnvClient()
//     if err != nil {
//         panic(err)
//     }
//
//     opts := types.ContainerStartOptions{}
//
//     containers, err := cli.ContainerList(context.Foreground(), "bitnami/redis:latest", opts )
//     if err != nil {
//         panic(err)
//     }
// }


func TestRedisQuicktimeStore( t *testing.T ) {
  red,err := CreateRedisQuicktimeStore("localhost:6379")
  if err != nil {
    t.Fatalf("Error creating Redis store: %s", err.Error() )
  }

src,err := lazyfs.OpenLocalFile( testfiles.TestMovPath )

  testKey := testfiles.TestMovPath
  lqt,err := red.Update( testKey, src )

  if err != nil {
    t.Errorf("Error reading quicktime: %s", err.Error() )
  }

  bar,has := red.Get( testfiles.TestMovPath )

  if !has {
    t.Errorf("Should have %s, but doesn't", testKey )
  //   } else if !reflect.DeepEqual( lqt, bar ) {
  //   t.Errorf("lqt and bar disagree")
  // }
  //
  // fmt.Println(lqt )
  //   fmt.Println( bar)

}
