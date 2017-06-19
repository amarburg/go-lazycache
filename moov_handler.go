package lazycache

import "fmt"
import "net/http"
import "strings"
import "strconv"

//import "strings"
import "io"
import "encoding/json"
import "github.com/amarburg/go-fast-png"
import "bytes"
import "regexp"
import "time"

import "github.com/amarburg/go-lazyquicktime"

var leadingNumbers, _ = regexp.Compile("^\\d+")

type QTMetadata struct {
	URL       string
	NumFrames int
	Duration  float32
}

//const qtPrefix = "qt."

var QTMetadataStore JSONStore

func init() {
	// Establish a default handler
	QTMetadataStore = CreateMapJSONStore()
}

func MoovHandler(node *Node, path []string, w http.ResponseWriter, req *http.Request) *Node {
	DefaultLogger.Log("msg", fmt.Sprintf("Quicktime handler: %s with residual path (%d): (%s)", node.Path, len(path), strings.Join(path, ":")))

	movStart := time.Now()

	// uri := node.Fs.Uri
	// uri.Path += node.Path
	//

	// Initialize or update as necessary
	lqt := &lazyquicktime.LazyQuicktime{}

	{
		DefaultLogger.Log("debug", "Locking metdatadata store")
		QTMetadataStore.Lock()
		defer QTMetadataStore.Unlock()

		DefaultLogger.Log("debug", "Querying metdatadata store")
		has, _ := QTMetadataStore.Get(node.trimPath, lqt)

		if !has {
			DefaultLogger.Log("debug", "Not in metdatadata store, querying CI")
			fs, err := node.Fs.LazyFile(node.Path)

			//fs, err := lazyfs.OpenHttpSource(uri)
			if err != nil {
				DefaultLogger.Log("err", "Something's went boom opening the HTTP Source!")
				http.Error(w, "Something's went boom opening the HTTP Source!", 500)
				return nil
			}

			DefaultLogger.Log("msg", fmt.Sprintf("Need to pull quicktime information for %s", fs.Path()))
			lqt, err = lazyquicktime.LoadMovMetadata(fs)
			if err != nil {
				DefaultLogger.Log("err", fmt.Sprintf("Something's went boom storing the quicktime file: %s", err.Error()))
				http.Error(w, fmt.Sprintf("Something's went boom storing the quicktime file: %s", err.Error()), 500)
				return nil
			}

			//fmt.Println(lqt)

			DefaultLogger.Log("msg", fmt.Sprintf("Updating metadata store for %s", fs.Path()))
			err = QTMetadataStore.Update(node.trimPath, *lqt)
			if err != nil {
				DefaultLogger.Log("err", fmt.Sprintf("Something's went boom storing the quicktime file: %s", err.Error()))
				http.Error(w, fmt.Sprintf("Something's went boom storing the quicktime file: %s", err.Error()), 500)
				return nil
			}
		} else {
			DefaultLogger.Log("msg", fmt.Sprintf("Map store had entry for %s", node.trimPath))
		}

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
			handleFrame(node, lqt, path[1:], w, req)
		default:
			http.Error(w, fmt.Sprintf("Didn't understand request \"%s\"", path[0]), 500)
		}
	}

	timeTrack(movStart, "Moov handler")

	return nil
}

func handleFrame(node *Node, lqt *lazyquicktime.LazyQuicktime, path []string, w http.ResponseWriter, req *http.Request) {

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

	UUID := req.URL.Path + ".png"
	url, ok := DefaultImageStore.Url(UUID)

	if ok {
		DefaultLogger.Log("msg", fmt.Sprintf("Image %s exists in the Image store at %s", UUID, url))
		// Set Content-Type or response
		w.Header().Set("Content-Type", "image/png")
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
		timeTrack(startExt, "Frame extraction")

		buffer := new(bytes.Buffer)

		encoder := fastpng.Encoder{
			CompressionLevel: fastpng.DefaultCompression,
		}

		startEncode := time.Now()
		err = encoder.Encode(buffer, img)
		timeTrack(startEncode, "Png encode")

		fmt.Println("PNG size", buffer.Len()/(1024*1024), "MB")

		imgReader := bytes.NewReader(buffer.Bytes())

		// write image to Image store
		DefaultImageStore.Store(UUID, imgReader)

		//Rewind the io, and write to the HTTP channel
		imgReader.Seek(0, io.SeekStart)
		_, err = imgReader.WriteTo(w)
		//fmt.Printf("Wrote %d bytes to http buffer\n", n)
		if err != nil {
			fmt.Printf("Error writing to HTTP buffer: %s\n", err.Error())
		}

		timeTrack(startExt, "Full extract and write")

	}

}
