package example

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/gorilla/mux"

	"github.com/charlieegan3/toolbelt/pkg/apis"
)

// JobsTool is an example tool which demonstrates the use of the toolbelt's job running functionality
type JobsTool struct {
	Count *int
}

func (jt *JobsTool) Name() string {
	return "jobs"
}

func (jt *JobsTool) FeatureSet() apis.FeatureSet {
	return apis.FeatureSet{
		Jobs: true,
	}
}

func (jt *JobsTool) Jobs() []apis.Job {
	return []apis.Job{&exampleJob{Count: jt.Count}}
}

func (jt *JobsTool) SetConfig(config map[string]any) error { return nil }
func (jt *JobsTool) DatabaseMigrations() (*embed.FS, string, error) {
	return nil, "", fmt.Errorf("not implemented")
}
func (jt *JobsTool) DatabaseSet(db *sql.DB)              {}
func (jt *JobsTool) HTTPPath() string                    { return "" }
func (jt *JobsTool) HTTPAttach(router *mux.Router) error { return nil }

// exampleJob shows a trivial apis.Job implementation
type exampleJob struct {
	Count *int
}

func (e *exampleJob) Name() string {
	return "example-job"
}

func (e *exampleJob) Run(ctx context.Context) error {
	doneCh := make(chan bool)
	errCh := make(chan error)

	go func() {
		*e.Count = *e.Count + 1
		fmt.Println(e.Name(), "ran")
		doneCh <- true
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-errCh:
		return fmt.Errorf("job failed with error: %s", e)
	case <-doneCh:
		return nil
	}
}

func (e *exampleJob) Timeout() time.Duration {
	return 3 * time.Second
}

func (e *exampleJob) Schedule() string {
	return "* * * * * *"
}
