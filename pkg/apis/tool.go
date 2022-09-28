package apis

import "github.com/gorilla/mux"

// FeatureSet is a list of optional features that a tool may use. Use of a feature will require more configuration from
// the tool belt.
type FeatureSet struct {
	// HTTP, if true, indicates that the tool needs to mount a subrouter on the tool belt webserver
	HTTP bool
	// Config, if true, indicates that the tool expects some config values to be passed at initialization time
	Config bool
	// Database, if true, indicates that the tool expects to connect to the tool belt database, owns a schema and will
	// have migrations that must be run.
	Database bool
}

type Tool interface {
	// Name returns the name of the tool for display purposes
	Name() string

	// FeatureSet returns the features that the tool uses and indicates to the tool belt how the tool should be
	// configured
	FeatureSet() FeatureSet

	// SetConfig sets the configuration for the tool
	SetConfig(config map[string]interface{}) error

	// HTTPPath returns the base path to use for the subrouter
	HTTPPath() string
	// HTTPAttach configures the tool's subrouter
	HTTPAttach(router *mux.Router) error
}
