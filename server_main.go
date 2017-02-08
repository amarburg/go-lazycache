package main

import (
	"flag"
	"fmt"
	"github.com/amarburg/go-lazycache/image_store"
	"github.com/amarburg/go-stoppable-http-server"
	"net/http"
	"net/url"
	"strings"
	//"github.com/amarburg/go-lazycache/quicktime_store"
)

var OOIRawDataRootURL = "https://rawdata.oceanobservatories.org/files/"

func main() {

	var (
		port          = flag.Int("port", 5000, "Network port")
		bind          = flag.String("bind", "127.0.0.1", "Network interface to bind")
		image_store   = flag.String("image_store", "", "")
		google_bucket = flag.String("image_store_bucket", "", "")
	)
	flag.Parse()

	//config,err := LoadLazyCacheConfig( *configFileFlag )

	// if err != nil {
	//   fmt.Printf("Error parsing config: %s\n", err.Error() )
	//   os.Exit(-1)
	// }

	//fmt.Println(config)
	ConfigureImageStore(*image_store, *google_bucket)

	server := StartLazycacheServer(*bind, *port)
	defer server.Stop()

	AddMirror(OOIRawDataRootURL)

	server.Wait()
}

func StartLazycacheServer(bind string, port int) *stoppable_http_server.SLServer {
	http.DefaultServeMux = http.NewServeMux()
	http.HandleFunc("/", IndexHandler)

	server := stoppable_http_server.StartServer(func(config *stoppable_http_server.HttpConfig) {
		config.Host = bind
		config.Port = port
	})

	return server
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
	root := fmt.Sprintf("/%s%s", strings.Join(splitHN, "/"), fs.Uri.Path)
	MakeRootNode(fs, root)

	RootMap[serverAddr] = root

	// fmt.Printf("Starting http handler at http://%s/\n", serverAddr)
	// fmt.Printf("Fs at http://%s%s\n", serverAddr, root )

	//http.ListenAndServe(serverAddr, nil)
}

func ConfigureImageStore(store_type string, bucket string) {
	switch strings.ToLower(store_type) {
	case "", "none":
		image_store.DefaultImageStore = image_store.NullImageStore{}
	case "google":
		fmt.Printf("Creating Google image store in bucket \"%s\"\n", bucket)
		image_store.DefaultImageStore = image_store.CreateGoogleStore(bucket)
	}
}
