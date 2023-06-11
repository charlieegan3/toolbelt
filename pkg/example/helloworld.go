package example

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/charlieegan3/toolbelt/pkg/apis"
	utilshttp "github.com/charlieegan3/toolbelt/pkg/utils/http"
)

// HelloWorld is an example tool which implements the Tool interface. It has a single http handler returning
// 'Hello World'
type HelloWorld struct{}

func (hw *HelloWorld) Name() string {
	return "hello-world"
}

func (hw *HelloWorld) FeatureSet() apis.FeatureSet {
	// This tool only uses the HTTP feature
	return apis.FeatureSet{
		HTTP: true,
	}
}

func (hw *HelloWorld) HTTPPath() string {
	return "example-hello-world"
}
func (hw *HelloWorld) HTTPHost() string {
	return ""
}

// SetConfig is a no-op for this tool
func (hw *HelloWorld) SetConfig(config map[string]any) error {
	return nil
}

func (hw *HelloWorld) HTTPAttach(router *mux.Router) error {
	router.HandleFunc("", utilshttp.BuildRedirectHandler(hw.HTTPPath()+"/")).Methods("GET")

	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte(hw.Name()))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	})

	return nil
}
