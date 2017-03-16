// +build integration

package lazycache

import "net/http"
import "fmt"
import "io/ioutil"

import "testing"

func TestHammerServer(t *testing.T) {
	server := StartLazycacheServer("127.0.0.1", 5000)
	defer server.Stop()
	AddMirror(OOIRawDataRootURL)

	HammerServer(1)
}

func HammerServer(count int) {
	c := make(chan int, count)

	for i := 0; i < count; i++ {
		go QueryServer(i, c)
	}

	for i := 0; i < count; i++ {
		<-c // wait for one task to complete
	}

}

func QueryServer(i int, c chan int) {
	resp, err := http.Get("https://127.0.0.1:5000/v1/org/oceanobservatories/rawdata/RS03ASHS/PN03B/06-CAMHDA301/2017/01/02/CAMHDA301-20170102T235100.mov")

	if err != nil {
		fmt.Printf("%d: ERROR: %v\n", i, err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%d: RESPONSE: %v\n%s\n", i, resp, body)
	}

	c <- i
}
