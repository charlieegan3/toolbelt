package apis

import (
	"database/sql"
	"embed"
	"github.com/gorilla/mux"
)

// FeatureSet is a list of optional features that a tool may use. Use of a feature will require more configuration from
// the tool belt.
type FeatureSet struct {
	// Config, if true, indicates that the tool expects some config values to be passed at initialization time
	Config bool
	// Database, if true, indicates that the tool expects to connect to the tool belt database, owns a schema and will
	// have migrations that must be run.
	Database bool

	// HTTP, if true, indicates that the tool needs to mount a subrouter on the tool belt webserver
	HTTP bool

	// HTTPHost, if true, indicates that the tool needs a subrouter with a host matcher
	HTTPHost bool

	// Jobs, if true, indicates that the tool has jobs which the belt must run
	Jobs bool

	// ExternalJobs, if true, indicates that the tool needs a function by which to start external jobs
	ExternalJobs bool
}

type Tool interface {
	// Name returns the name of the tool for display purposes
	Name() string

	// FeatureSet returns the features that the tool uses and indicates to the tool belt how the tool should be
	// configured
	FeatureSet() FeatureSet

	// SetConfig sets the configuration for the tool
	SetConfig(config map[string]interface{}) error

	// DatabaseMigrate runs the database migrations for the tool
	DatabaseMigrations() (*embed.FS, string, error)
	// DatabaseSet sets the database connection for the tool
	DatabaseSet(db *sql.DB)

	// HTTPPath returns the base path to use for the subrouter
	HTTPPath() string
	// HTTPHost returns the host to use for the subrouter, if not blank
	HTTPHost() string
	// HTTPAttach configures the tool's subrouter
	HTTPAttach(router *mux.Router) error

	// Jobs returns a list of jobs that the tool defines and needs to have run
	Jobs() ([]Job, error)

	// ExternalJobsFunc sets the function that the tool can use to start external jobs
	ExternalJobsFuncSet(func(job ExternalJob) error)
}
