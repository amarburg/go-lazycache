package main

import (
       "fmt"
       "github.com/amarburg/go-stoppable-http-server"
       "net/http"
       "net/url"
       "strings"
)

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
		 DefaultImageStore = NullImageStore{}
	case "google":
	   DefaultImageStore = CreateGoogleStore(bucket)
	}
}
