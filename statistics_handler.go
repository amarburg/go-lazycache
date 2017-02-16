package lazycache


import "fmt"
import "net/http"
import "encoding/json"
import "time"


var startTime = time.Now()


func StatisticsHandler(w http.ResponseWriter, req *http.Request) {

  stats := struct{
    uptime     time.Duration
    UptimeSec  float64

    ImageStore  interface{}
    DirectoryStore interface{}
    QuicktimeStore interface{}

    Roots       map[string]interface{}
  }{
  Roots:   make( map[string]interface{} ),
}

  // Populate
  stats.uptime = time.Since( startTime )
  stats.UptimeSec = stats.uptime.Seconds()

  stats.ImageStore = DefaultImageStore.Statistics()
  stats.DirectoryStore = listing_store.Statistics()
  stats.QuicktimeStore = quicktime_store.Statistics()

  for root,_ := range RootMap {
    stats.Roots[root] = struct{
      Fs  interface{}   `json: "fs"`
      }{
      Fs: RootMap[root].node.Fs.Statistics,
    }
  }

		b, err := json.MarshalIndent(stats, "", "  ")
		if err != nil {
			fmt.Fprintln(w, "JSON error:", err)
		}

		w.Write(b)
}
