package main

import "net/http"
import "fmt"
import "io/ioutil"



func main() {
	HammerServer(10)
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
	resp, err := http.Get("http://localhost:8080/v1/org/oceanobservatories/rawdata/files/RS03ASHS/PN03B/06-CAMHDA301/2016/03/09/CAMHDA301-20160309T090000Z.mov/frame/10000")

	if err != nil {
		fmt.Printf("%d: ERROR: %v\n", i, err)
	} else {
		defer resp.Body.Close()
		ioutil.ReadAll(resp.Body)
		//fmt.Printf("%d: RESPONSE: %v\n%s\n", i, resp, body)
	}

	c <- i
}
