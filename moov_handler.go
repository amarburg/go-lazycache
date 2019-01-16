package lazycache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/amarburg/go-fast-png"
	"github.com/spf13/viper"
	"golang.org/x/image/bmp"
 	"github.com/disintegration/imaging"
	"image/jpeg"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

import "github.com/amarburg/go-lazyquicktime"

var leadingNumbers, _ = regexp.Compile("^\\d+")

//go:generate $GOPATH/bin/easyjson -all $GOFILE

// MoovHandlerTiming
type MoovHandlerTiming struct {
	// Don't export start times, so they don't get JSON-encoded
	// Todo, I could clean up this API a bit...
	handlerStart                          time.Time
	Handler, Metadata, Extraction, Encode time.Duration
}

type moovOutputMetadata struct {
	URL       string
	NumFrames uint64
	Duration  float64
	FileSize  int64
}

type QTEntry struct {
	lqt      *lazyquicktime.LazyQuicktime
	metadata moovOutputMetadata
}

type QTStore struct {
	Cache map[string](*QTEntry)
	Mutex sync.Mutex

	Stats struct {
		Requests, Misses int64
	}
}

var qtCache QTStore

func init() {
	qtCache = QTStore{
		Cache: make(map[string](*QTEntry)),
	}
}

func (cache *QTStore) getLQT(node *Node) (*QTEntry, error) {

	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	cache.Stats.Requests++

	// Initialize or update as necessary
	Logger.Log("debug", fmt.Sprintf("Querying metadata store for %s", node.Path))
	qte, has := cache.Cache[node.trimPath]

	if !has {
		cache.Stats.Misses++

		//Logger.Log("msg", fmt.Errorf("Initializing LazyFile to %s", node.Path))
		fs, err := node.Fs.FileSource(node.Path)

		if err != nil {
			return nil, fmt.Errorf("Something went boom while opening the HTTP Source!")
		}

		//block, err := lazyfs.OpenBlockStore( fs, 20 )
		// if err != nil {
		// 	return nil, fmt.Errorf("Something went boom while opening the HTTP Source!")
		// }

		//Logger.Log("msg", fmt.Sprintf("Need to pull quicktime information for %s", fs.Path()))
		lqt, err := lazyquicktime.LoadMovMetadata(fs)
		if err != nil {
			return nil, fmt.Errorf("Something went boom while storing the quicktime file: %s", err.Error())
		}

		//Logger.Log("msg", fmt.Sprintf("Updating metadata store for %s", fs.Path()))
		qte = &QTEntry{
			lqt: lqt,
			metadata: moovOutputMetadata{
				FileSize:  lqt.FileSize,
				URL:       node.Path,
				NumFrames: lqt.NumFrames(),
				Duration:  lqt.Duration().Seconds(),
			},
		}

		cache.Cache[node.trimPath] = qte

	} else {
		Logger.Log("msg", fmt.Sprintf("lqt cache has entry for %s", node.Path))
	}

	return qte, nil
}

func MoovHandler(node *Node, path []string, w http.ResponseWriter, req *http.Request) *Node {
	Logger.Log("msg", fmt.Sprintf("Quicktime handler: %s with residual path (%d): (%s)", node.Path, len(path), strings.Join(path, ":")))

	timing := MoovHandlerTiming{
		handlerStart: time.Now(),
	}

	// uri := node.Fs.Uri
	// uri.Path += node.Path

	metadataStart := time.Now()
	qte, err := qtCache.getLQT(node)
	timeTrack(metadataStart, &timing.Metadata)

	if err != nil {
		Logger.Log("msg", err.Error())

		b, _ := json.MarshalIndent(struct {
			URL, Error string
		}{
			URL:   node.Path,
			Error: err.Error(),
		}, "", "  ")

		w.Write(b)
		return nil
	}

	if len(path) == 0 {
		// Leaf node
		startEncode := time.Now()

		Logger.Log("msg", fmt.Sprintf("Returning movie information for %s", node.Path))

		b, err := qte.metadata.MarshalJSON()
		if err != nil {
			fmt.Fprintln(w, "JSON error:", err)
		}

		timeTrack(startEncode, &timing.Encode)
		timeTrack(timing.handlerStart, &timing.Handler)

		w.Header().Set("X-lazycache-timing-handler-ns", strconv.Itoa(int(timing.Handler.Nanoseconds())))
		w.Header().Set("X-lazycache-timing-metadata-ns", strconv.Itoa(int(timing.Metadata.Nanoseconds())))

		w.Write(b)
	} else {

		// Handle any residual path elements (frames, etc) here
		switch strings.ToLower(path[0]) {
		case "frame":
			extractFrame(node, qte, path[1:], w, req, &timing)
		default:
			http.Error(w, fmt.Sprintf("Didn't understand request \"%s\"", path[0]), 500)
		}
	}

	t, _ := timing.MarshalJSON()
	Logger.Log("timing", t)

	return nil
}

func extractFrame(node *Node, qte *QTEntry, path []string, w http.ResponseWriter, req *http.Request, timing *MoovHandlerTiming) {

	if qte == nil || qte.lqt == nil {
		http.Error(w, "Error in extractFrame", 500)
		return
	}

	if len(path) == 0 {
		http.Error(w, fmt.Sprintf("Need to specify frame number"), 500)
		return

	}

	frameNum, err := strconv.Atoi(leadingNumbers.FindString(path[0]))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing frame number \"%s\"", path[0]), 500)
		return
	}

	if uint64(frameNum) > qte.metadata.NumFrames {
		http.Error(w, fmt.Sprintf("Requested frame %d in movie of length %d frames", frameNum, qte.metadata.NumFrames), 400)
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
	case ".bmp":
		contentType = "image/bmp"
		extension = ".bmp"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
		extension = ".jpg"
	case "", ".png":
		extension = ".png"
		contentType = "image/png"
	case ".rgba", ".raw":
		extension = ".rgba"
		contentType = "image/x-raw-rgba"
	default:
		http.Error(w, fmt.Sprintf("Unknown image extension \"%s\"", extension), 500)
		return
	}

	UUID := req.URL.Path + extension
	url, ok := ImageCache.Url(UUID)

	if ok {
		Logger.Log("msg", fmt.Sprintf("Image %s exists in the Image store at %s", UUID, url))
		// Set Content-Type or response
		w.Header().Set("Content-Type", contentType)
		Logger.Log("msg", fmt.Sprintf("Redirecting to %s", url))
		http.Redirect(w, req, url, http.StatusTemporaryRedirect)

	} else {

		startExt := time.Now()
		img, err := qte.lqt.ExtractNRGBA(uint64(frameNum))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error generating image for frame %d: %s", frameNum, err.Error()), 500)
			return
		}
		timeTrack(startExt, &timing.Extraction)

		query := req.URL.Query()

		widthStr, widthValid := query["width"]
		heightStr, heightValid := query["height"]


		if( widthValid && heightValid ) {

			width, _ := strconv.Atoi(widthStr[0])
			height, _ := strconv.Atoi(heightStr[0])

			Logger.Log("msg", fmt.Sprintf("Resizing to %d x %d", width, height))
			resized := imaging.Resize(img, width, height, imaging.Lanczos)
			img = resized

		} else if( widthValid || heightValid ) {

			http.Error(w, fmt.Sprintf("Both width and height must be specified.  Got width = %s and height = %s ", widthStr, heightStr), 500)
			return

		} else {

			// Check HTTP header for scale information
			hdrWidth := req.Header.Get("X-lazycache-output-width")
			hdrHeight := req.Header.Get("X-lazycache-output-height")

			if( len(hdrWidth) > 0 && len(hdrHeight) > 0 ) {

				width, _ := strconv.Atoi(hdrWidth)
				height, _ := strconv.Atoi(hdrHeight)

				Logger.Log("msg", fmt.Sprintf("Resizing to %d x %d", width, height))
				resized := imaging.Resize(img, width, height, imaging.Lanczos)
				img = resized
			}
		}

		startEncode := time.Now()

		var imgReader *bytes.Reader

		switch contentType {
		case "image/png":

			buffer := new(bytes.Buffer)
			encoder := new(fastpng.Encoder)

			// TODO, allow configuration of PNGs
			if viper.GetBool("public") {
				encoder.CompressionLevel = fastpng.BestCompression
			}

			// {
			// 	CompressionLevel: fastpng.BestSpeed,
			// }

			err = encoder.Encode(buffer, img)
			imgReader = bytes.NewReader(buffer.Bytes())

		case "image/jpeg":
			buffer := new(bytes.Buffer)
			err = jpeg.Encode(buffer, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
			imgReader = bytes.NewReader(buffer.Bytes())

		case "image/bmp":
			buffer := new(bytes.Buffer)
			err = bmp.Encode(buffer, img)
			imgReader = bytes.NewReader(buffer.Bytes())

		case "image/x-raw-rgba":
			if viper.GetBool("allow-raw-output") {
				// stand-in
				//buffer = img.Pix
				imgReader = bytes.NewReader(img.Pix)
			} else {
				http.Error(w, "This server is not configured to produce raw output.", 501)
				return
			}
		}

		timeTrack(startEncode, &timing.Encode)

		//Logger.Log("debug", fmt.Sprintf("%s size %d MB\n", contentType, buffer.Len()/(1024*1024)))

		// write image to Image store
		ImageCache.Store(UUID, imgReader)

		timeTrack(timing.handlerStart, &timing.Handler)

		// Add timing information to HTTP Header
		w.Header().Set("X-lazycache-timing-handler-ns", strconv.Itoa(int(timing.Handler.Nanoseconds())))
		w.Header().Set("X-lazycache-timing-metadata-ns", strconv.Itoa(int(timing.Metadata.Nanoseconds())))
		w.Header().Set("X-lazycache-timing-extraction-ns", strconv.Itoa(int(timing.Extraction.Nanoseconds())))
		w.Header().Set("X-lazycache-timing-encode-ns", strconv.Itoa(int(timing.Encode.Nanoseconds())))

		imgReader.Seek(0, io.SeekStart)
		_, err = imgReader.WriteTo(w)
		if err != nil {
			fmt.Printf("Error writing to HTTP buffer: %s\n", err.Error())
		}

	}

}
