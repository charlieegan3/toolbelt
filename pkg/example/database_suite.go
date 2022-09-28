package example

import (
	"database/sql"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/charlieegan3/toolbelt/pkg/tool"
)

type ExampleDatabaseToolSuite struct {
	suite.Suite
	DB *sql.DB
}

func (s ExampleDatabaseToolSuite) Run(t *testing.T) {
	suite.Run(t, &s)
}

func (s *ExampleDatabaseToolSuite) TestDatabaseTool() {
	t := s.T()

	tb := tool.NewBelt()

	tb.SetDatabase(s.DB)

	databaseTool := &DatabaseTool{}

	err := tb.AddTool(databaseTool)
	require.NoError(t, err)

	c := tb.StartServer("0.0.0.0", "9032")
	require.NoError(t, err)
	defer func() {
		c <- os.Interrupt
		err = tb.StopServer(5 * time.Second)
		require.NoError(t, err)
	}()

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
