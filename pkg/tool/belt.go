package tool

import (
	"fmt"
	"github.com/gorilla/mux"

	"github.com/charlieegan3/toolbelt/pkg/apis"
)

// Belt is the main struct for the Tool Belt. It contains the base router which all tools are registered to
type Belt struct {
	Router *mux.Router
}

// NewBelt creates a new Belt struct with an initalized router
func NewBelt() Belt {
	return Belt{
		Router: mux.NewRouter(),
	}
}

// AddTool adds a new tool to the belt. Each tool is given a subrouter with the base path set to the tool's HTTPPath
func (b Belt) AddTool(tool apis.Tool) error {
	path := tool.HTTPPath()
	toolRouter := b.Router.PathPrefix(fmt.Sprintf("/%s", path)).Subrouter()
	err := tool.HTTPAttach(toolRouter)
	if err != nil {
		return fmt.Errorf("failed to attach tool: %v", err)
	}

	return nil
}
