// build +profile
package main

import "testing"
import "net/http"
import "encoding/json"
import "image/png"

func jsonQuery(t *testing.T, url string, result interface{}) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		t.Errorf("Got error response from local server")
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = decoder.Decode(result)

	if err != nil {
		t.Errorf("Got error decoding JSON RootMap: %s", err.Error())
	}
}

func TestRoot(t *testing.T) {
	server := RunOOIServer("127.0.0.1", 12345)
	defer server.Stop()

	rootmap := make(map[string]interface{})
	jsonQuery(t, "http://127.0.0.1:12345/", &rootmap)

	if len(rootmap) == 0 {
		t.Error("Zero-length RootMap")
	}
}

func TestOOIRoot(t *testing.T) {
	server := RunOOIServer("127.0.0.1", 12345)
	defer server.Stop()

	rootmap := make(map[string]interface{})
	jsonQuery(t, "http://127.0.0.1:12345/v1/org/oceanobservatories/rawdata/files/", &rootmap)

	if rootmap["Files"] == nil {
		t.Errorf("Path=\"Files\" doesn't exist")
	}

	if rootmap["Path"] != "/" {
		t.Errorf("Expected Path=\"/\", got: %s", rootmap["Path"])
	}
}

func TestOOIRootMovieMetadata(t *testing.T) {
	server := RunOOIServer("127.0.0.1", 12345)
	defer server.Stop()

	metadata := make(map[string]interface{})
	jsonQuery(t, "http://127.0.0.1:12345/v1/org/oceanobservatories/rawdata/files/RS03ASHS/PN03B/06-CAMHDA301/2016/07/24/CAMHDA301-20160724T030000Z.mov", &metadata)

	// Movie length known apriori
	if metadata["NumFrames"].(float64) != 25162 {
		t.Error("Incorrect metadata")
	}

	//TODO: image checks
}

func TestOOIRootMovieMetadataParallel(t *testing.T) {
	server := RunOOIServer("127.0.0.1", 12345)
	defer server.Stop()

	results := make(chan float64)

	numClients := 100
	for i := 0; i < numClients; i++ {
		go func(out chan<- float64) {
			metadata := make(map[string]interface{})
			jsonQuery(t, "http://127.0.0.1:12345/v1/org/oceanobservatories/rawdata/files/RS03ASHS/PN03B/06-CAMHDA301/2016/07/24/CAMHDA301-20160724T030000Z.mov", &metadata)

			out <- metadata["NumFrames"].(float64)
		}(results)
	}

	for i := 0; i < numClients; i++ {
		numFrames := <-results

		// Movie length known apriori
		if numFrames != 25162.0 {
			t.Error("Incorrect metadata")
		}
	}

	//TODO: image checks
}

func TestOOIRootImageDecode(t *testing.T) {
	server := RunOOIServer("127.0.0.1", 12345)
	defer server.Stop()

	resp, err := http.Get("http://127.0.0.1:12345/v1/org/oceanobservatories/rawdata/files/RS03ASHS/PN03B/06-CAMHDA301/2016/07/24/CAMHDA301-20160724T030000Z.mov/frame/1000")
	//resp,err := http.Get("http://127.0.0.1:12345/v1/localhost:9081/RS03ASHS/PN03B/06-CAMHDA301/2016/07/24/CAMHDA301-20160724T030000Z.mov/frame/1000")
	defer resp.Body.Close()

	if err != nil {
		t.Errorf("Error retrieving image: %s", err.Error())
	}

	img, err := png.Decode(resp.Body)

	if err != nil {
		t.Errorf("Error decoding image: %s", err.Error())
	}

	if img.Bounds().Dx() != 1920 || img.Bounds().Dy() != 1080 {
		t.Errorf("Extracted image not expected size (was %d x %d)", img.Bounds().Dx(), img.Bounds().Dy())
	}

	// Other image checks
}
