package lazycache

import (
	//"fmt"
	"github.com/amarburg/go-lazyfs"
	testfiles "github.com/amarburg/go-lazyfs-testfiles"
	"github.com/amarburg/go-lazyquicktime"
	"testing"
	//"reflect"
)

func TestMapJsonStore(t *testing.T) {
	red := CreateMapJSONStore()

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
	// fmt.Println(lqt)
	// fmt.Println(retrieved)

}
