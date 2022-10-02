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

// Jobs is an example tool which demonstrates the use of the toolbelt's job running functionality
type Jobs struct{}

func (hw *Jobs) Name() string {
	return "jobs"
}

func (hw *Jobs) FeatureSet() apis.FeatureSet {
	return apis.FeatureSet{
		Jobs: true,
	}
}

func (hw *Jobs) Jobs() []apis.Job {
	return []apis.Job{&exampleJob{}}
}

func (hw *Jobs) SetConfig(config map[string]any) error { return nil }
func (hw *Jobs) DatabaseMigrations() (*embed.FS, string, error) {
	return nil, "", fmt.Errorf("not implemented")
}
func (hw *Jobs) DatabaseSet(db *sql.DB)              {}
func (hw *Jobs) HTTPPath() string                    { return "" }
func (hw *Jobs) HTTPAttach(router *mux.Router) error { return nil }

// exampleJob shows a trivial apis.Job implementation
type exampleJob struct{}

func (e *exampleJob) Name() string {
	return "example-job"
}

func (e *exampleJob) Run(ctx context.Context) error {
	doneCh := make(chan bool)
	errCh := make(chan error)

	go func() {
		time.Sleep(1 * time.Second)
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
	return "*/10 * * * * *"
}
