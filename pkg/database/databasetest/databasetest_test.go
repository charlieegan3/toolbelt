package databasetest

import (
	"github.com/charlieegan3/toolbelt/pkg/example"
	"testing"

	_ "github.com/lib/pq"
)

func TestDatabaseSuite(t *testing.T) {
	s := &DatabaseSuite{
		ConfigPath: "../../../config.test.yaml",
	}

	s.Setup(t)

	s.AddDependentSuite(&example.ExampleDatabaseToolSuite{DB: s.DB})
	s.AddDependentSuite(&example.ExampleJobsToolSuite{DB: s.DB})

	s.Run(t)
}
