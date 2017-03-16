package lazycache

import "fmt"
import "net/http"
import "encoding/json"

func IndexHandler(w http.ResponseWriter, req *http.Request) {

	// Map RootMap to a different structure

	type RootMapOut struct {
		APIPathV1 string
	}

	jsonRootMap := make(map[string]RootMapOut)

	for key, root := range RootMap {
		jsonRootMap[key] = RootMapOut{
			APIPathV1: root.node.trimPath,
		}
	}

	//	if jsonExtension.MatchString(req.URL.Path) {

	b, err := json.MarshalIndent(jsonRootMap, "", "  ")
	if err != nil {
		fmt.Fprintln(w, "JSON error:", err)
	}

	w.Write(b)
	//	} else {

	// 	fmt.Fprintf(w, "<html><body><ul>")
	// 	for key, val := range RootMap {
	// 		fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>\n", val, key)
	// 	}
	// 	fmt.Fprintf(w, "</ul></body></html>")
	// 	//fmt.Println("Indexing from ", req.URL.String())
	// }
}
