package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/charlieegan3/toolbelt/pkg/example"
	"github.com/charlieegan3/toolbelt/pkg/tool"
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

	count := 0
	err = tb.AddTool(&example.JobsTool{Count: &count})
	if err != nil {
		log.Fatalf("failed to add tool: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-c:
			cancel()
		}
	}()

	go tb.RunJobs(ctx)

	tb.RunServer(ctx, "0.0.0.0", "3000")
}
