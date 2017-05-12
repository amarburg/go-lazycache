package lazycache

import "net/http"

// Cache-Control: no-cache, no-store, must-revalidate
// Pragma: no-cache
// Expires: 0
func addCacheDefeatHeaders(w http.ResponseWriter) {
	w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}
