package lazycache

import "fmt"
import "net/http"
import "strings"
import "strconv"
import "path/filepath"

//import "strings"
import "io"
import "encoding/json"
import "github.com/amarburg/go-fast-png"
import "image/jpeg"
import "bytes"
import "regexp"
import "time"

import "sync"

import "github.com/amarburg/go-lazyquicktime"

var leadingNumbers, _ = regexp.Compile("^\\d+")

type QTMetadata struct {
	URL       string
	NumFrames int
	Duration  float32
}

type QTStore struct {
	Cache map[string](*lazyquicktime.LazyQuicktime)
	Mutex sync.Mutex

	Stats struct {
		Requests, Misses int64
	}
}

var qtCache QTStore

func init() {
	qtCache = QTStore{
		Cache: make(map[string](*lazyquicktime.LazyQuicktime)),
	}
}

type MoovHandlerTiming struct {
	Handler, Extraction, Encode time.Duration
}

func retrieveLazyQuicktime(node *Node) (*lazyquicktime.LazyQuicktime, error) {

	//DefaultLogger.Log("debug", fmt.Sprintf("Locking metadata store for %s", node.Path))
	qtCache.Mutex.Lock()
	defer qtCache.Mutex.Unlock()

	qtCache.Stats.Requests++

	// Initialize or update as necessary
	//DefaultLogger.Log("debug", fmt.Sprintf("Querying metadata store for %s", node.Path))
	lqt, has := qtCache.Cache[node.trimPath]

	if !has {
		qtCache.Stats.Misses++
		fs, err := node.Fs.LazyFile(node.Path)

		if err != nil {
			return nil, fmt.Errorf("Something's went boom opening the HTTP Source!")
		}

		//DefaultLogger.Log("msg", fmt.Sprintf("Need to pull quicktime information for %s", fs.Path()))
		lqt, err = lazyquicktime.LoadMovMetadata(fs)
		if err != nil {
			return nil, fmt.Errorf("Something's went boom storing the quicktime file: %s", err.Error())
		}

		//DefaultLogger.Log("msg", fmt.Sprintf("Updating metadata store for %s", fs.Path()))
		qtCache.Cache[node.trimPath] = lqt
		if err != nil {
			return nil, fmt.Errorf("Something's went boom storing the quicktime file: %s", err.Error())
		}

	}
	//else {
	// Could record cache misses
	// DefaultLogger.Log("msg", fmt.Sprintf("Map store had entry for %s", node.trimPath))
	//}

	return lqt, nil
}

func MoovHandler(node *Node, path []string, w http.ResponseWriter, req *http.Request) *Node {
	DefaultLogger.Log("msg", fmt.Sprintf("Quicktime handler: %s with residual path (%d): (%s)", node.Path, len(path), strings.Join(path, ":")))

	timing := MoovHandlerTiming{}
	movStart := time.Now()

	// uri := node.Fs.Uri
	// uri.Path += node.Path
	//
	lqt, err := retrieveLazyQuicktime(node)

	if err != nil {
		DefaultLogger.Log("msg", err.Error())
		http.Error(w, err.Error(), 500)
	}

	if len(path) == 0 {
		// Leaf node

		DefaultLogger.Log("msg", fmt.Sprintf("Returning movie information for %s", node.Path))

		out := QTMetadata{
			URL:       node.Path,
			NumFrames: lqt.NumFrames(),
			Duration:  lqt.Duration(),
		}

		// Temporary structure for JSON output
		b, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			fmt.Fprintln(w, "JSON error:", err)
		}

		DefaultLogger.Log("msg", fmt.Sprintf("..... done"))

		w.Write(b)
	} else {

		// Handle any residual path elements (frames, etc) here
		switch strings.ToLower(path[0]) {
		case "frame":
			extractFrame(node, lqt, path[1:], w, req, &timing)
		default:
			http.Error(w, fmt.Sprintf("Didn't understand request \"%s\"", path[0]), 500)
		}
	}

	timeTrack(movStart, &timing.Handler)

	t, _ := json.Marshal(timing)
	DefaultLogger.Log("timing", t)

	return nil
}

func extractFrame(node *Node, lqt *lazyquicktime.LazyQuicktime, path []string, w http.ResponseWriter, req *http.Request, timing *MoovHandlerTiming) {

	if len(path) == 0 {
		http.Error(w, fmt.Sprintf("Need to specify frame number"), 500)
		return

	}

	frameNum, err := strconv.Atoi(leadingNumbers.FindString(path[0]))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing frame number \"%s\"", path[0]), 500)
		return
	}

	if frameNum > lqt.NumFrames() {
		http.Error(w, fmt.Sprintf("Requested frame %d in movie of length %d frames", frameNum, lqt.NumFrames()), 400)
		return
	}

	if frameNum < 1 {
		http.Error(w, "Requested frame 0, Quicktime movies start with frame 1", 400)
		return
	}

	// Looks for extension
	extension := filepath.Ext(path[0])

	var contentType string

	switch extension {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
		extension = ".jpg"
	case "", ".png":
		extension = ".png"
		contentType = "image/png"
	default:
		http.Error(w, fmt.Sprintf("Unknown image extension \"%s\"", extension), 500)
		return
	}

	UUID := req.URL.Path + extension
	url, ok := ImageCache.Url(UUID)

	if ok {
		DefaultLogger.Log("msg", fmt.Sprintf("Image %s exists in the Image store at %s", UUID, url))
		// Set Content-Type or response
		w.Header().Set("Content-Type", contentType)
		// w.Header().Set("Location", url)
		DefaultLogger.Log("msg", fmt.Sprintf("Redirecting to %s", url))
		http.Redirect(w, req, url, http.StatusTemporaryRedirect)

	} else {

		startExt := time.Now()
		img, err := lqt.ExtractFrame(frameNum)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error generating image for frame %d: %s", frameNum, err.Error()), 500)
			return
		}
		timeTrack(startExt, &timing.Extraction)

		buffer := new(bytes.Buffer)

		startEncode := time.Now()

		switch contentType {
		case "image/png":
			encoder := fastpng.Encoder{
				CompressionLevel: fastpng.DefaultCompression,
			}
			err = encoder.Encode(buffer, img)
		case "image/jpeg":
			err = jpeg.Encode(buffer, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		}

		timeTrack(startEncode, &timing.Encode)

		DefaultLogger.Log("debug", fmt.Sprintf("%s size %d MB\n", contentType, buffer.Len()/(1024*1024)))

		imgReader := bytes.NewReader(buffer.Bytes())

		// write image to Image store
		ImageCache.Store(UUID, imgReader)

		//Rewind the io, and write to the HTTP channel
		imgReader.Seek(0, io.SeekStart)
		_, err = imgReader.WriteTo(w)
		//fmt.Printf("Wrote %d bytes to http buffer\n", n)
		if err != nil {
			fmt.Printf("Error writing to HTTP buffer: %s\n", err.Error())
		}

		// timeTrack(startExt, "Full extract and write")

	}

}
