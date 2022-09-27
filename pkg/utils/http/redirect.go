package http

import "net/http"

// BuildRedirectHandler returns a handler that redirects to the given path
func BuildRedirectHandler(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusSeeOther)
	}
}
