package example

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/charlieegan3/toolbelt/pkg/tool"
)

type ExampleDatabaseToolSuite struct {
	suite.Suite
	DB *sql.DB
}

func (s *ExampleDatabaseToolSuite) Run(t *testing.T) {
	suite.Run(t, s)
}

func (s *ExampleDatabaseToolSuite) TestDatabaseTool() {
	t := s.T()

	tb := tool.NewBelt()

	tb.SetDatabase(s.DB)

	databaseTool := &DatabaseTool{}

	err := tb.AddTool(context.Background(), databaseTool)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go tb.RunServer(ctx, "0.0.0.0", "9032")

	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost:9032",
			Path:   "/database/",
		},
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	if resp.StatusCode != http.StatusOK {
		t.Log(string(body))
		t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if string(body) != "database value" {
		t.Errorf("expected 'database value', got '%s'", string(body))
	}

	err = tb.DatabaseDownMigrate(databaseTool)
	require.NoError(t, err)
}
