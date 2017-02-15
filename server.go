package lazycache

import (
       "fmt"
       "net/url"
       "strings"
       "net/http"
       kitlog "github.com/go-kit/kit/log"
)

const ApiVersion = "v1"

func RegisterDefaultHandlers() {
  http.HandleFunc("/v1/statistics/", StatisticsHandler )
  http.HandleFunc("/", IndexHandler)
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


	RootMap[serverAddr] = root
}

func ConfigureImageStore(store_type string, bucket string, logger kitlog.Logger) {
	switch strings.ToLower(store_type) {
	case "", "none":
		 DefaultImageStore = NullImageStore{}
	case "google":
	   DefaultImageStore = CreateGoogleStore(bucket, logger)
	}
}
