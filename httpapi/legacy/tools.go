package legacy

import (
	"net/http"
	"fmt"
)

// Store form parameter value of specified key to *val and return true; if key does not exist, set HTTP status 400 and return false.
func Require(w http.ResponseWriter, r *http.Request, key string, val *string) bool {
	*val = r.FormValue(key)
	if *val == "" {
		http.Error(w, fmt.Sprintf("Please pass POST/PUT/GET parameter value of '%s'.", key), 400)
		return false
	}
	return true
}