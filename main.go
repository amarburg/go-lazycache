package main

import (
	"flag"
)

var OOIRawDataRootURL = "https://rawdata.oceanobservatories.org/files/"

func main() {

	var (
		port          = flag.Int("port", 5000, "Network port to listen on (default: 5000)")
		bind          = flag.String("bind", "0.0.0.0", "Network interface to bind to (defaults to 0.0.0.0)")
		image_store   = flag.String("image-store", "", "Type of image store (none, google)")
		google_bucket = flag.String("image-store-bucket", "", "Bucket used for Google image store")
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
