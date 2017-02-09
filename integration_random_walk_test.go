// +build integration

package main

import (
 "github.com/amarburg/go-lazycache/stress_test"
 "testing"
)



func TestRandomWalk(t *testing.T) {
	flag.Parse()

	server := StartLazycacheServer("127.0.0.1", 5000)
	defer server.Stop()

	AddMirror(OOIRawDataRootURL)

	err := stress_test.RandomWalk("127.0.0.1:5000")
	if err != nil {
		t.Error(err)
	}
}
