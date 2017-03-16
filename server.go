package lazycache

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	//       kitlog "github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const ApiVersion = "v1"

func RegisterDefaultHandlers() {
	http.HandleFunc("/", IndexHandler)
	http.Handle("/metrics/", promhttp.Handler())
}

func AddMirror(serverAddr string) {

	url, err := url.Parse(serverAddr)
	fs, err := OpenHttpFS(*url)

	if err != nil {
		panic(fmt.Sprintf("Error opening HTTP FS Source: %s", err.Error()))
	}

	//serverAddr := fmt.Sprintf("%s:%d", config.ServerIp, config.ServerPort)

	// Reverse hostname
	splitHN := MungeHostname(fs.Uri.Host)
	root := fmt.Sprintf("/%s/%s%s", ApiVersion, strings.Join(splitHN, "/"), fs.Uri.Path)
	MakeRootNode(fs, root)
}
