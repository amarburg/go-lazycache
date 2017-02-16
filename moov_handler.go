package lazycache

import "fmt"
import "net/http"
import "strings"
import "strconv"

//import "strings"
import "io"
import "encoding/json"
import "image/png"
import "bytes"

import "github.com/amarburg/go-lazyfs"

import "github.com/amarburg/go-lazyquicktime"

func MoovHandler(node *Node, path []string, w http.ResponseWriter, req *http.Request) *Node {
	//  fmt.Fprintf( w, "Quicktime handler: %s with residual path (%d): (%s)\n", node.Path, len(path), strings.Join(path,":") )

	uri := node.Fs.Uri
	uri.Path += node.Path

	// Initialize or update as necessary
	lqt, ok := Get(node.trimPath)
	if !ok {
		node.updateMutex.Lock()
		lqt, ok = Get(node.trimPath)
		if !ok {
			fs, err := lazyfs.OpenHttpSource(uri)
			if err != nil {
				http.Error(w, "Something's went boom opening the HTTP Soruce!", 500)
				return nil
			}

			lqt, err = Update(node.trimPath, fs)
			if err != nil {
				http.Error(w, fmt.Sprintf("Something's went boom parsing the quicktime file: %s", err.Error()), 500)
				return nil
			}
		}
		node.updateMutex.Unlock()
	}

	if len(path) == 0 {
		// Leaf node

		// Temporary structure for JSON output
		out := struct {
			URL       string
			NumFrames int
			Duration  float32
		}{
			URL:       uri.String(),
			NumFrames: lqt.NumFrames(),
			Duration:  lqt.Duration(),
		}

		b, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			fmt.Fprintln(w, "JSON error:", err)
		}

		w.Write(b)
	} else {

		// Handle any residual path elements (frames, etc) here
		switch strings.ToLower(path[0]) {
		case "frame":
			handleFrame(node, lqt, path[1:], w, req)
		default:
			http.Error(w, fmt.Sprintf("Didn't understand request \"%s\"", path[0]), 500)
		}
	}

	return nil
}

func handleFrame(node *Node, lqt *lazyquicktime.LazyQuicktime, path []string, w http.ResponseWriter, req *http.Request) {

	if len(path) == 0 {
		http.Error(w, fmt.Sprintf("Need to specify frame number"), 500)
		return

	}

	frameNum, err := strconv.Atoi(path[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing frame number \"%s\"", path[0]), 500)
		return
	}

	if frameNum > lqt.NumFrames() {
		http.Error(w, fmt.Sprintf("Requested frame %d in movie with %d frames", frameNum, lqt.NumFrames()), 400)
		return
	}
	if frameNum < 1 {
		http.Error(w, "Requested frame 0, Quicktime movies start with frame 1", 400)
		return
	}

	UUID := req.URL.Path + ".png"

	url, ok := DefaultImageStore.Url(UUID)

	if ok {
		fmt.Printf("Image %s exists in the Image store at %s", UUID, url)
		// Set Content-Type or response
		w.Header().Set("Content-Type", "image/png")
		http.Redirect(w, req, url, 307)

	} else {

		img, err := lqt.ExtractFrame(frameNum)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error generating image for frame %d: %s", frameNum, err.Error()), 500)
			return
		}

		buffer := new(bytes.Buffer)
		err = png.Encode(buffer, img)

		imgReader := bytes.NewReader(buffer.Bytes())

		// write image to Image store
		DefaultImageStore.Store(UUID, imgReader)

		//fmt.Println(buffer)
		imgReader.Seek(0, io.SeekStart)

		// Write image to HTTP request
		n, err := imgReader.WriteTo(w)
		fmt.Printf("Wrote %d bytes to http buffer\n", n)
		if err != nil {
			fmt.Printf("Error writing to HTTP buffer: %s\n", err.Error())
		}

	}
}
