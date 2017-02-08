// +build integration

package main

import "net/http"
import "fmt"
import "encoding/json"

import "github.com/amarburg/go-lazycache/listing_store"
import "math/rand"
import "flag"

import "testing"


var randomWalkCount = flag.Int("random-walk-count", 100, "Number of random walk queries to make")
var randomWalkParallelism = flag.Int("random-walk-parallelism", 1, "Parallelism of random walk queries")



func TestRandomWalk(t *testing.T) {
	server := StartLazycacheServer("127.0.0.1", 5000)
	defer server.Stop()

	AddMirror(OOIRawDataRootURL)

	RandomWalk( *randomWalkCount, *randomWalkParallelism )
}

func RandomWalk(count, parallelism int) {

	urls := make( chan string, parallelism )
	out  := make( chan bool, count )

	urls <- "http://localhost:5000/org/oceanobservatories/rawdata/files/"

	for i := 0; i < count; i++ {
		go RandomWalkQuery( urls, out )
	}

	for i := 0; i < count; i++ {
		fmt.Printf("Wait for task to complete %d",i)
		resp := <-out // wait for one task to complete

		if !resp {
			close(urls)
			return
		}
	}

}

func RandomWalkQuery( urls chan string, out chan bool ) {
	url,ok := <- urls

	if !ok {
		out <- false
		fmt.Println("Channel closed")

		return
	}

	fmt.Println("Random walker Querying URL", url)

	resp, err := http.Get( url )
	if err != nil {
		fmt.Printf("%d: ERROR: %s\n", url, err)
	} else {
		defer resp.Body.Close()
		// body, _ := ioutil.ReadAll(resp.Body)
		// fmt.Printf("%d: RESPONSE: %v\n%s\n", i, resp, body)

		// Parse response
		decoder := json.NewDecoder(resp.Body )
		var listing listing_store.DirListing
		err := decoder.Decode( listing )

		if err == nil {
			fmt.Printf("Has %d directories\n", len(listing.Directories) )


			a := rand.Intn( len(listing.Directories ) )
			b := rand.Intn( len(listing.Directories ) )

			urls <- url + listing.Directories[a]
			urls <- url + listing.Directories[b]

		}

		fmt.Println("Good response")
		out <- true
		return
	}

	fmt.Println("Bad response")
	out <- false
}
