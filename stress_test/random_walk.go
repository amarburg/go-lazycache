package stress_test

import (
  "flag"
  "fmt"
  "errors"
  "net/http"
  "encoding/json"
  "github.com/amarburg/go-lazycache/listing_store"
  "math/rand"
)

var RandomWalkCount = flag.Int("random-walk-count", 50, "Number of random walk queries to make")
var RandomWalkParallelism = flag.Int("random-walk-parallelism", 5, "Parallelism of random walk queries")

func RandomWalk( host string ) error {

  var urls = make(chan string, *RandomWalkCount*2 )
  var out = make(chan bool)

  count := *RandomWalkCount
  parallelism := *RandomWalkParallelism

  baseurl := fmt.Sprintf("http://%s/org/oceanobservatories/rawdata/files/", host )


	for i := 0; i < parallelism; i++ {
		go RandomWalkQuery(urls,out, baseurl)
    urls <- baseurl 
	}


	//urls <- fmt.Sprintf("http://%s/org/oceanobservatories/rawdata/files/RS03ASHS/PN03B/06-CAMHDA301/", host )

	i := 0
	for {
		fmt.Printf("Wait for task %d to complete ...", i)
		resp := <-out // wait for one task to complete

		// Always seed the channel with another url, just in case
		//urls <- fmt.Sprintf("http://%s/org/oceanobservatories/rawdata/files/",host)

		i++

		if !resp {
			return errors.New("Error from child")
		} else if i > count {
			return nil
		}
	}

}

func RandomWalkQuery(urls chan string, out chan bool, baseurl string) {
	fmt.Println("In random walker")
	for url := range urls {

		fmt.Println("Random walker Querying URL", url)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%d: ERROR: %s\n", url, err)
			fmt.Printf("Error making request: %s\n", err.Error())
			out <- false
			return

		}

		defer resp.Body.Close()
		// body, _ := ioutil.ReadAll(resp.Body)
		// fmt.Printf("%d: RESPONSE: %v\n%s\n", i, resp, body)

		// Parse response
		decoder := json.NewDecoder(resp.Body)
		var listing listing_store.DirListing

		if err := decoder.Decode(&listing); err != nil {
			fmt.Println("Error reading response: %s\n", err.Error())
			out <- false
			return
		}
		//fmt.Printf("Has %d directories\n", len(listing.Directories))

		if len(listing.Directories) > 0 {

			urls <- url + listing.Directories[rand.Intn(len(listing.Directories))]
			//urls <- url + listing.Directories[rand.Intn(len(listing.Directories))]
		} else {
      urls <- baseurl
    }

		//fmt.Println("Good response")
		out <- true
	}
}
