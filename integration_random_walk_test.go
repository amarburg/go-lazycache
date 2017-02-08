// +build integration

package main

import "net/http"
import "fmt"
import "encoding/json"

import "github.com/amarburg/go-lazycache/listing_store"
import "math/rand"
import "flag"
import "errors"

import "testing"

var randomWalkCount = flag.Int("random-walk-count", 10, "Number of random walk queries to make")
var randomWalkParallelism = flag.Int("random-walk-parallelism", 5, "Parallelism of random walk queries")

func TestRandomWalk(t *testing.T) {
	server := StartLazycacheServer("127.0.0.1", 5000)
	defer server.Stop()

	AddMirror(OOIRawDataRootURL)

	err := RandomWalk(*randomWalkCount, *randomWalkParallelism)
	if err != nil {
		t.Error(err)
	}
}

var urls = make(chan string, *randomWalkCount)
var out = make(chan bool)

func RandomWalk(count, parallelism int) error {

	for i := 0; i < parallelism; i++ {
		go RandomWalkQuery()
	}

	urls <- "http://localhost:5000/org/oceanobservatories/rawdata/files/"

	i := 0
	for {
		fmt.Printf("Wait for task to complete %d ...", i)
		resp := <-out // wait for one task to complete
		i++
		fmt.Printf("Done\n")

		if !resp {
			return errors.New("Error from child")
		} else if i > count {
			return nil
		}
	}

}

func RandomWalkQuery() {
	fmt.Println("In random walker")
	for url := range urls {

		fmt.Println("Random walker Querying URL", url)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%d: ERROR: %s\n", url, err)
			fmt.Printf("Error making request: %s\n", err.Error())
			out <- false
		} else {
			defer resp.Body.Close()
			// body, _ := ioutil.ReadAll(resp.Body)
			// fmt.Printf("%d: RESPONSE: %v\n%s\n", i, resp, body)

			// Parse response
			decoder := json.NewDecoder(resp.Body)
			var listing listing_store.DirListing

			if err := decoder.Decode( &listing ); err == nil {
				fmt.Printf("Has %d directories\n", len(listing.Directories))

				if len(listing.Directories) > 0 {

				urls <- url + listing.Directories[rand.Intn(len(listing.Directories))]
				urls <- url + listing.Directories[rand.Intn(len(listing.Directories))]
			}

				fmt.Println("Good response")
				out <- true
			} else {
				fmt.Println("Error reading response: %s\n", err.Error())
				out <- false
			}
		}

	}
}
