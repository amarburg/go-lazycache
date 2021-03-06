package lazycache

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

func ViperConfiguration() {

	// Configuration
	viper.SetDefault("port", 8080)
	viper.SetDefault("bind", "0.0.0.0")
	viper.SetDefault("imagestore", "")
	viper.SetDefault("imagestore.bucket", "camhd-image-cache")

	viper.SetDefault("quicktimestore", "")
	viper.SetDefault("directorystore", "")
	viper.SetDefault("redishost", "localhost:6379")

	viper.SetDefault("fileoverlay", "")
	viper.SetDefault("fileoverlay.flatten", "")

	viper.SetConfigName("lazycache")
	viper.AddConfigPath("/etc/lazycache")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	flag.Int("port", 80, "Network port to listen on (default: 8080)")
	flag.String("bind", "0.0.0.0", "Network interface to bind to (defaults to 0.0.0.0)")

	flag.String("image-store", "", "Type of image store (none, local, google)")
	flag.String("image-store-bucket", "camhd-image-cache", "Bucket used for Google image store")
	flag.String("image-store-root", "", "Path to local image store directory (must be writable)")
	flag.String("image-store-url", "", "Root URL for webserver which serves image store directory")

	flag.String("file-overlay", "", "Path to local file overlay")
	viper.BindPFlag("fileoverlay", flag.Lookup("file-overlay"))

	flag.Bool("file-overlay-flatten", false, "Do flatten the file overlay")
	viper.BindPFlag("fileoverlay.flatten", flag.Lookup("file-overlay-flatten"))

	// flag.String("quicktime-store", "", "Type of quicktime store (none, redis)")
	// flag.String("directory-store", "", "Type of directory store (none, redis)")
	flag.String("redis-host", "localhost:6379", "Host used for redis store")

	viper.BindPFlag("port", flag.Lookup("port"))
	viper.BindPFlag("bind", flag.Lookup("bind"))

	viper.BindPFlag("imagestore", flag.Lookup("image-store"))
	viper.BindPFlag("imagestore.bucket", flag.Lookup("image-store-bucket"))
	viper.BindPFlag("imagestore.root", flag.Lookup("image-store-root"))
	viper.BindPFlag("imagestore.url", flag.Lookup("image-store-url"))

	//viper.BindPFlag("directorystore", flag.Lookup("directory-store"))
	//viper.BindPFlag("quicktimestore", flag.Lookup("quicktime-store"))
	viper.BindPFlag("redishost", flag.Lookup("redis-host"))

	flag.Bool("allow-raw-output", false, "Allow images to be output as raw bytestrings using PIL.Image.tobytes()")
	viper.BindPFlag("allow-raw-output", flag.Lookup("allow-raw-output"))

	flag.Bool("public", false, "Set public mode")
	viper.BindPFlag("public", flag.Lookup("public"))

	flag.Parse()
}

func ConfigureImageStoreFromViper() {
	storeKey := viper.GetString("imagestore")
	Logger.Log("msg", fmt.Sprintf("Configuring image store with type \"%s\"", storeKey))
	switch strings.ToLower(storeKey) {
	default:
		Logger.Log("msg", fmt.Sprintf("Unable to determine type of image store from \"%s\"", storeKey))
		ImageCache = NullImageStore{}
	case "", "none":
		Logger.Log("msg", "No image store configured.")
		ImageCache = NullImageStore{}
	case "local":
		ImageCache = CreateLocalStore(viper.GetString("imagestore.root"), viper.GetString("imagestore.url"))
	case "google":
		ImageCache = CreateGoogleStore(viper.GetString("imagestore.bucket"))
	}
}

func ConfigureFromViper() {
	ViperConfiguration()

	Logger.Log("msg", "In ConfigureFromViper")
	ConfigureImageStoreFromViper()

	if viper.GetBool("allow-raw-output") {
		Logger.Log("msg", "Raw image output enabled.")
	}

	if viper.GetBool("public") {
		Logger.Log("msg", "Enabled public server mode.")
	}
}
