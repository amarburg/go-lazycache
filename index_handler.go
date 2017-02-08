package main

import "fmt"
import "net/http"
import "encoding/json"

var RootMap = make(map[string]string)

func IndexHandler(w http.ResponseWriter, req *http.Request) {

	b, err := json.MarshalIndent(RootMap, "", "  ")
	if err != nil {
		fmt.Fprintln(w, "JSON error:", err)
	}

	w.Write(b)

	// fmt.Fprintf(w, "<html><body><ul>")
	// for key,val := range RootMap {
	// fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>\n", val, key)
	// }
	// fmt.Fprintf(w, "</ul></body></html>")
	//fmt.Println("Indexing from ", req.URL.String())
}
