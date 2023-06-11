package databasetest

import (
	"testing"

	_ "github.com/lib/pq"

	"github.com/charlieegan3/toolbelt/pkg/example"
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
