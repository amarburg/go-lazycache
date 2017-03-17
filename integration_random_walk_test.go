// +build integration

package lazycache

import (
	stress "github.com/amarburg/go-lazycache-benchmarking"
	"github.com/amarburg/go-stoppable-http-server"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

var OOIRawDataRootURL = "https://rawdata.oceanobservatories.org/files/"

func StartLazycacheServer(bind string, port int) *stoppable_http_server.SLServer {
	http.DefaultServeMux = http.NewServeMux()
	http.HandleFunc("/", IndexHandler)

	server := stoppable_http_server.StartServer(func(config *stoppable_http_server.HttpConfig) {
		config.Host = bind
		config.Port = port
	})

	return server
}

func TestRandomWalk(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	server := StartLazycacheServer("127.0.0.1", 5000)
	defer server.Stop()

	AddMirror(OOIRawDataRootURL)

	settings := stress.NewSettings()

	settings.SetCount(100)
	settings.SetParallelism(5)
	err := stress.RandomWalk(*settings, "http://127.0.0.1:5000/v1/org/oceanobservatories/rawdata/files/")

	if err != nil {
		t.Error(err)
	}
}
