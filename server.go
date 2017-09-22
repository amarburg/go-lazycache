package lazycache

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	//       kitlog "github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/amarburg/go-lazyfs"
	"github.com/amarburg/go-lazyquicktime"
	"github.com/amarburg/go-quicktime"

	"github.com/spf13/viper"
)

const ApiVersion = "v1"

const Version = "v0.1.0"

func RegisterDefaultHandlers() {
	http.HandleFunc("/", RootHandler)
	http.Handle("/metrics/", promhttp.Handler())
	http.HandleFunc("/info/", InfoHandler)
}

func AddMirror(serverAddr string) {

	url, err := url.Parse(serverAddr)
	fs, err := OpenHttpFS(*url)

	if err != nil {
		panic(fmt.Sprintf("Error opening HTTP FS Source: %s", err.Error()))
	}

	// Reverse hostname
	splitHN := MungeHostname(fs.Uri.Host)
	root := fmt.Sprintf("/%s/%s%s", ApiVersion, strings.Join(splitHN, "/"), fs.Uri.Path)

	if len(viper.GetString("fileoverlay")) > 0 {
		ofs, err := OpenFileOverlayFS(fs, viper.GetString("fileoverlay"))

		if viper.GetBool("fileoverlay.flatten") {
			ofs.Flatten = true
		}

		Logger.Log("msg", fmt.Sprintf("Adding fileoverlay at %s, flatten=%v",
			viper.GetString("fileoverlay"), ofs.Flatten))

		if err != nil {
			panic(fmt.Sprintf("Error opening FileOverlay FS Source: %s", err.Error()))
		}

		MakeRootNode(ofs, root)

	} else {
		MakeRootNode(fs, root)
	}
}

func InfoHandler(w http.ResponseWriter, req *http.Request) {
	b := &bytes.Buffer{}

	// Could this be automatic...?

	fmt.Fprintf(b, "{\n")
	fmt.Fprintf(b, "    \"amarburg/github/go-lazycache\": \"%s\",\n", Version)
	fmt.Fprintf(b, "    \"amarburg/github/go-lazyquicktime\": \"%s\",\n", lazyquicktime.Version)
	fmt.Fprintf(b, "    \"amarburg/github/go-lazyfs\": \"%s\",\n", lazyfs.Version)
	fmt.Fprintf(b, "    \"amarburg/github/go-quicktime\": \"%s\",\n", quicktime.Version)

	fmt.Fprintf(b, "}\n")

	w.Write(b.Bytes())
}
