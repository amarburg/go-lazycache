package main

import (
	"fmt"
	stress "github.com/amarburg/go-lazycache-benchmarking"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	settings := stress.NewSettings()
	settings.SetCount(100)
	settings.SetParallelism(5)
	err := stress.RandomWalk(*settings, "http://localhost:8080/v1/org/oceanobservatories/rawdata/files/")

	if err != nil {
		fmt.Sprintf("Error!   %s", err)
	}
}
