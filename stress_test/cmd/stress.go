package main


import (
  "github.com/amarburg/go-lazycache/stress_test"
  "flag"
)

var host = flag.String("host","127.0.0.1:5000","Host to access")


func main() {

  flag.Parse()

  err := stress_test.RandomWalk( *host)

	if err != nil {
		panic(err)
	}

}
