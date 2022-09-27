package apis

import "github.com/gorilla/mux"

type Tool interface {
	// Name returns the name of the tool for display purposes
	Name() string

	// HTTPPath returns the base path to use for the subrouter
	HTTPPath() string
	// HTTPAttach configures the tool's subrouter
	HTTPAttach(router *mux.Router) error
}
