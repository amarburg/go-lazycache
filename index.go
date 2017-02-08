package main

import "fmt"
import "net/http"

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<a href=\"org/oceanobservatories/rawdata//\">rawdata.oceanobservatories.org/</a>\n")
	fmt.Println("Indexing from ", req.URL.String())
}
