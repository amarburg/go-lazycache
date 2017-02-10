// +build integration

package main

import (
 stress "github.com/amarburg/go-lazycache-benchmarking"
 "testing"
 "math/rand"
 "time"
)

func TestRandomWalk(t *testing.T) {
	rand.Seed( time.Now().UTC().UnixNano())

	server := StartLazycacheServer("127.0.0.1", 5000)
	defer server.Stop()

	AddMirror(OOIRawDataRootURL)

	err := stress.RandomWalk( stress.AddUrl("http://127.0.0.1:5000/org/oceanobservatories/rawdata/files/"),
 													stress.SetCount( 100 ),
													stress.SetParallelism( 5 ) )
	if err != nil {
		t.Error(err)
	}
}
