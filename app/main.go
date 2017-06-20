package main

import (
	"fmt"
	"github.com/amarburg/go-lazycache"
	"github.com/amarburg/go-stoppable-http-server"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
)

var ooiRawDataRootURL = "https://rawdata.oceanobservatories.org/files/"

// StartLazycacheServer starts a HTTP server at http://bind:port/ and registers default Lazycache handlers
func StartLazycacheServer(bind string, port int) *stoppable_http_server.SLServer {
	http.DefaultServeMux = http.NewServeMux()

	msg := fmt.Sprintf("Listening on http://%s:%d/", bind, port)
	lazycache.DefaultLogger.Log("msg", msg)

	server := stoppable_http_server.StartServer(func(config *stoppable_http_server.HttpConfig) {
		config.Host = bind
		config.Port = port
	})

	lazycache.RegisterDefaultHandlers()

	return server
}

// RunOOIServer starts an Lazycache server and registers the standard
// Rutgers rawdata destination
func RunOOIServer(bind string, port int) *stoppable_http_server.SLServer {
	server := StartLazycacheServer(bind, port)

	lazycache.AddMirror(ooiRawDataRootURL)

	return server
}

func main() {

	// Add my own options
	viper.SetDefault("cpuprofile", "")
	flag.String("cpuprofile", "", "CPU Profile file")
	viper.BindPFlag("cpuprofile", flag.Lookup("cpuprofile"))

	lazycache.ConfigureFromViper()

	if viper.GetString("cpuprofile") != "" {
		fmt.Println("Creating cpu profile \"", viper.GetString("cpuprofile"), "\"")
		f, err := os.Create(viper.GetString("cpuprofile"))
		if err != nil {
			log.Fatal(err)
		}

		if err = pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}

		defer f.Close()
		defer pprof.StopCPUProfile()
	}

	server := RunOOIServer(viper.GetString("bind"), viper.GetInt("port"))
	defer server.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			server.Stop()
		}
	}()

	// Handle Google App Engine health requests
	http.HandleFunc("/_ah/health", healthCheckHandler)

	server.Wait()

}
