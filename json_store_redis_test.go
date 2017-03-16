package lazycache

import (
	"testing"
	//"context"
	// "github.com/docker/docker/api/types"
	// "github.com/docker/docker/client"
	"fmt"
	"github.com/amarburg/go-lazyfs"
	testfiles "github.com/amarburg/go-lazyfs-testfiles"
	"github.com/amarburg/go-lazyquicktime"
	//"reflect"
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

func TestRedixJsonStore(t *testing.T) {
	red, err := CreateRedisJSONStore("localhost:6379", "test")
	if err != nil {
		t.Fatalf("Error creating Redis store: %s", err.Error())
	}

	src, err := lazyfs.OpenLocalFile(testfiles.TestMovPath)
	lqt, err := lazyquicktime.LoadMovMetadata(src)

	testKey := testfiles.TestMovPath

	err = red.Update(testKey, lqt)

	if err != nil {
		t.Errorf("Error reading quicktime: %s", err.Error())
	}

	var retrieved lazyquicktime.LazyQuicktime
	ok, err := red.Get(testfiles.TestMovPath, &retrieved)

	if err != nil {
		t.Errorf("Got error when retrieving value: %s", err.Error())
	} else if !ok {
		t.Errorf("Should have %s, but doesn't", testKey)
	}
	//   } else if !reflect.DeepEqual( lqt, bar ) {
	//   t.Errorf("lqt and bar disagree")
	//}
	//
	fmt.Println(lqt)
	fmt.Println(retrieved)

}
