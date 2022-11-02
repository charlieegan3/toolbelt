package example

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"

	"github.com/charlieegan3/toolbelt/pkg/tool"
)

type ExampleJobsToolSuite struct {
	suite.Suite
	DB *sql.DB
}

func (s *ExampleJobsToolSuite) Run(t *testing.T) {
	suite.Run(t, s)
}

func (s *ExampleJobsToolSuite) TestJobsTool() {
	t := s.T()

	tb := tool.NewBelt()

	tb.SetDatabase(s.DB)

	var count int

	jobsTool := &JobsTool{Count: &count}

	err := tb.AddTool(context.Background(), jobsTool)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go tb.RunJobs(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		time.Sleep(time.Second)
		wg.Done()
	}()

	wg.Wait()

	require.Greaterf(t, count, 0, "example job should have run at least once")
}
