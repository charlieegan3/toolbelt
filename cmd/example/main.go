package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/charlieegan3/toolbelt/pkg/example"
	"github.com/charlieegan3/toolbelt/pkg/tool"
)

// this is an example use of a Tool Belt showing the registration of an example Hello World tool
func main() {
	tb := tool.NewBelt()

	err := tb.AddTool(example.HelloWorld{})
	if err != nil {
		log.Fatalf("failed to add tool: %v", err)
	}

	srv := &http.Server{
		Handler:      tb.Router,
		Addr:         fmt.Sprintf("%s:%s", "0.0.0.0", "3000"),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
