package lazycache

import (
       "fmt"
       "net/url"
       "strings"
       kitlog "github.com/go-kit/kit/log"
)



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

func ConfigureImageStore(store_type string, bucket string, logger kitlog.Logger) {
	switch strings.ToLower(store_type) {
	case "", "none":
		 DefaultImageStore = NullImageStore{}
	case "google":
	   DefaultImageStore = CreateGoogleStore(bucket, logger)
	}
}
