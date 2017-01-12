package lazycache

import "fmt"
import "net/http"

func MoovHandler( w http.ResponseWriter, req *http.Request ) {
  http.Error( w, fmt.Sprintf("Error serving file %s", req.URL.Path ), 404 )
}
