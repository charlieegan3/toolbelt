package main

import (
	"github.com/charlieegan3/toolbelt/pkg/example"
	"github.com/charlieegan3/toolbelt/pkg/tool"
	"log"
	"time"
)

// this is an example use of a Tool Belt showing the registration of an example Hello World tool
func main() {
	tb := tool.NewBelt()

	// this might be loaded from disk in some real example
	tb.SetConfig(map[string]any{
		"config-tool": map[string]any{
			"exampleValue": "example config value",
		},
	})

	err := tb.AddTool(&example.HelloWorld{})
	if err != nil {
		log.Fatalf("failed to add tool: %v", err)
	}

	err = tb.AddTool(&example.ConfigTool{})
	if err != nil {
		log.Fatalf("failed to add tool: %v", err)
	}

	err = tb.AddTool(&example.Jobs{})
	if err != nil {
		log.Fatalf("failed to add tool: %v", err)
	}

	tb.RunJobs()

	c := tb.StartServer("0.0.0.0", "3000")

	<-c // wait for interrupt

	err = tb.StopServer(5 * time.Second)
	if err != nil {
		log.Fatalf("failed to stop server: %v", err)
	}
}
