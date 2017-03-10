package lazycache

import (
	"github.com/spf13/viper"
	flag "github.com/spf13/pflag"
	"fmt"
	"strings"
)

func ViperConfiguration() {

	// Configuration
	viper.SetDefault("port", 8080 )
	viper.SetDefault("bind","0.0.0.0")
	viper.SetDefault("imagestore", "")
	viper.SetDefault("imagestore.bucket","camhd-image-cache")

	viper.SetConfigName("lazycache")
	viper.AddConfigPath("/etc/lazycache")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil  { // Handle errors reading the config file
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// ignore
		default:
	    panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	viper.SetEnvPrefix("lazycache")
	viper.AutomaticEnv()
	// Convert '.' to '_' in configuration variable names
	viper.SetEnvKeyReplacer( strings.NewReplacer(".", "_") )

	// var (
	// 	bindFlag          = flag.String("bind", "0.0.0.0", "Network interface to bind to (defaults to 0.0.0.0)")
	// 	ImageStoreFlag   = flag.String("image-store", "", "Type of image store (none, google)")
	// 	ImageBucketFlag = flag.String("image-store-bucket", "", "Bucket used for Google image store")
	// )
	flag.Int("port", 80, "Network port to listen on (default: 8080)")
	flag.String("bind", "0.0.0.0", "Network interface to bind to (defaults to 0.0.0.0)")
	flag.String("image-store", "", "Type of image store (none, google)")

	flag.String("image-store-bucket", "camhd-image-cache", "Bucket used for Google image store")
	flag.String("image-local-root", "", "Bucket used for Google image store")
	flag.String("image-url-root", "", "Bucket used for Google image store")


	viper.BindPFlag("port", flag.Lookup("port"))
	viper.BindPFlag("bind", flag.Lookup("bind"))
	viper.BindPFlag("imagestore", flag.Lookup("image-store"))

	viper.BindPFlag("imagestore.bucket", flag.Lookup("image-store-bucket"))
	viper.BindPFlag("imagestore.localroot", flag.Lookup("image-local-root"))
	viper.BindPFlag("imagestore.urlroot", flag.Lookup("image-url-root"))


	flag.Parse()
}

func ConfigureImageStoreFromViper() {
	switch strings.ToLower( viper.GetString("imagestore" )) {
	case "", "none":
		fmt.Printf("Unable to determine type of image store from \"%s\"", viper.GetString("imagestore" ) )
		 DefaultImageStore = NullImageStore{}
	case "local":
			DefaultImageStore = CreateLocalStore(viper.GetString("imagestore.localRoot"),
																						viper.GetString("imagestore.urlRoot") )
	case "google":
	   DefaultImageStore = CreateGoogleStore(viper.GetString("imagestore.bucket") )
	}

}
