package example

import (
	"github.com/gorilla/mux"
	"net/http"

	utilshttp "github.com/charlieegan3/toolbelt/pkg/utils/http"
)

// HelloWorld is an example tool which implements the Tool interface. It has a single http handler returning
// 'Hello World'
type HelloWorld struct{}

func (hw HelloWorld) Name() string {
	return "Hello World"
}

func (hw HelloWorld) HTTPPath() string {
	return "hello-world"
}

func (hw HelloWorld) HTTPAttach(router *mux.Router) error {
	router.HandleFunc("", utilshttp.BuildRedirectHandler(hw.HTTPPath()+"/")).Methods("GET")

	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte(hw.Name()))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	})

	return nil
}
