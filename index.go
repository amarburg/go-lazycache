package lazycache

import "fmt"
import "net/http"

func Index(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(w, "<a href=\"rawdata.oceanobservatories.org/\">rawdata.oceanobservatories.org/</a>\n")
  fmt.Println("Indexing from ", req.URL.String() )
}
