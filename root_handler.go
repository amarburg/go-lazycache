package lazycache

import "fmt"
import "net/http"
import "encoding/json"

// RootHandler is the default HTTP handler, registered at "/"
// It returns a JSON structure giving the relative path to
// each of the registered mirrors.
//
// e.g.
//
//   {
//     "https://rawdata.oceanobservatories.org/files/": {
//       "APIPath": {
//         "V1": "/v1/org/oceanobservatories/rawdata/files/"
//       }
//     }
//   }
func RootHandler(w http.ResponseWriter, req *http.Request) {

	// Temporary structures which define the output JSON structure
	type APIPathOut struct {
		V1 string
	}

	type RootMapOut struct {
		APIPath APIPathOut
	}

	jsonRootMap := make(map[string]RootMapOut)

	for key, root := range RootMap {
		jsonRootMap[key] = RootMapOut{
			APIPath: APIPathOut{
				V1: root.node.trimPath,
			},
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
