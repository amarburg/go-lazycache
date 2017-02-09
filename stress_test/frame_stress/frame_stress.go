package main

import (
	"flag"
	"github.com/amarburg/go-lazycache/stress_test"
)

var imgUrl = flag.String("url", "", "URL of img to query")

func main() {

	flag.Parse()

	err := stress_test.FrameStress(*imgUrl)

	if err != nil {
		panic(err)
	}

}
