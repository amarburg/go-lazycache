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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// var (
	// 	bindFlag          = flag.String("bind", "0.0.0.0", "Network interface to bind to (defaults to 0.0.0.0)")
	// 	ImageStoreFlag   = flag.String("image-store", "", "Type of image store (none, google)")
	// 	ImageBucketFlag = flag.String("image-store-bucket", "", "Bucket used for Google image store")
	// )
	flag.Int("port", 80, "Network port to listen on (default: 8080)")
	flag.String("bind", "0.0.0.0", "Network interface to bind to (defaults to 0.0.0.0)")

	flag.String("image-store", "", "Type of image store (none, local, google)")
	flag.String("image-store-bucket", "camhd-image-cache", "Bucket used for Google image store")
	flag.String("image-local-root", "", "Bucket used for Google image store")
	flag.String("image-url-root", "", "Bucket used for Google image store")

	flag.String("quicktime-store", "", "Type of quicktime store (none, redis)")
	flag.String("directory-store", "", "Type of directory store (none, redis)")
	flag.String("redis-host", "localhost:6379", "Host used for redis store")

	viper.BindPFlag("port", flag.Lookup("port"))
	viper.BindPFlag("bind", flag.Lookup("bind"))
	viper.BindPFlag("imagestore", flag.Lookup("image-store"))

	viper.BindPFlag("imagestore.bucket", flag.Lookup("image-store-bucket"))
	viper.BindPFlag("imagestore.localroot", flag.Lookup("image-local-root"))

	viper.BindPFlag("directorystore", flag.Lookup("directory-store"))
	viper.BindPFlag("quicktimestore", flag.Lookup("quicktime-store"))
	viper.BindPFlag("redishost", flag.Lookup("redis-host"))

	flag.Parse()
}

func ConfigureImageStoreFromViper() {
	switch strings.ToLower(viper.GetString("imagestore")) {
	default:
		DefaultLogger.Log("msg", "Unable to determine type of image store from \"%s\"", viper.GetString("imagestore"))
		DefaultImageStore = NullImageStore{}
	case "", "none":
		DefaultLogger.Log("msg", "No image store configured.")
		DefaultImageStore = NullImageStore{}
	case "local":
		DefaultImageStore = CreateLocalStore(viper.GetString("imagestore.localRoot"),
			viper.GetString("imagestore.bind"))
	case "google":
		DefaultImageStore = CreateGoogleStore(viper.GetString("imagestore.bucket"))
	}
}

func ConfigureQuicktimeStoreFromViper() {

	switch strings.ToLower(viper.GetString("quicktimestore")) {
	default:
		DefaultLogger.Log("msg", "Unable to determine type of image store from \"%s\"", viper.GetString("quicktimestore"))
		QTMetadataStore = CreateMapJSONStore()
	case "", "none":
		DefaultLogger.Log("msg", "Using default QuicktimeStore.")
		QTMetadataStore = CreateMapJSONStore()
	case "redis":
		hostname := viper.GetString("redishost")
		redis, err := CreateRedisJSONStore(hostname, "qt")
		if err != nil {
			DefaultLogger.Log("msg", fmt.Sprintf("Failed to configure Redis Quicktime store to host \"%s\"", hostname))
		}

		DefaultLogger.Log("msg", fmt.Sprintf("Logging movie metadata to Redis at %s", hostname))
		QTMetadataStore = redis
	}
}

func ConfigureDirectoryStoreFromViper() {

	switch strings.ToLower(viper.GetString("directorystore")) {
	default:
		DefaultLogger.Log("msg", "Unable to determine type of directory store from \"%s\"", viper.GetString("directorystore"))
		DirKeyStore = CreateMapJSONStore()
	case "", "none":
		DefaultLogger.Log("msg", "Using default directory store.")
		DirKeyStore = CreateMapJSONStore()
	case "redis":
		hostname := viper.GetString("redishost")
		redis, err := CreateRedisJSONStore(hostname, "dir")
		if err != nil {
			DefaultLogger.Log("msg", fmt.Sprintf("Failed to configure Redis directory store to host \"%s\"", hostname))
		}

		DefaultLogger.Log("msg", fmt.Sprintf("Logging directory metadata to Redis at %s", hostname))
		DirKeyStore = redis
	}
}

func ConfigureFromViper() {
	ConfigureImageStoreFromViper()
	ConfigureDirectoryStoreFromViper()
	ConfigureQuicktimeStoreFromViper()
}
