package databasetest

import (
	"database/sql"
	"github.com/charlieegan3/toolbelt/pkg/database"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DependentSuite interface {
	Run(t *testing.T)
}

// DatabaseSuite is the top of the test suite hierarchy for all tests that use
// the database.
type DatabaseSuite struct {
	ConfigPath string
	DB         *sql.DB

	suite.Suite
	suites []DependentSuite
}

func (s *DatabaseSuite) Setup(t *testing.T) {
	defaultConfigPath := "../../../config.test.yaml"
	if s.ConfigPath == "" {
		s.ConfigPath = defaultConfigPath
	}

	viper.SetConfigFile(s.ConfigPath)
	err := viper.ReadInConfig()
	if err != nil {
		t.Fatalf("failed to load test config: %s", err)
	}

	// initialize a database connection to init the db
	params := viper.GetStringMapString("database.params")
	connectionString := viper.GetString("database.connectionString")
	db, err := database.Init(connectionString, params, "postgres", false)
	if err != nil {
		t.Fatalf("failed to init DB: %s", err)
	}

	// dbname must be set to a test db name
	dbname, ok := params["dbname"]
	if !ok {
		t.Fatalf("test dbname param was not set, failing as unsure what DB to use: %s", err)
	}

	// if the database exists, then we drop it to give a clean test state
	// this happens at the start of the test suite so that the state is there
	// after a test run to inspect if need be
	exists, err := database.Exists(db, dbname)
	if err != nil {
		t.Fatalf("failed to check if test DB exists: %s", err)
	}
	if exists {
		// drop existing test db
		err = database.Drop(db, dbname)
		if err != nil {
			t.Fatalf("failed to drop test database: %s", err)
		}
	}

	// create the test db for this test run
	err = database.Create(db, dbname)
	if err != nil {
		t.Fatalf("failed to create test database: %s", err)
	}

	// init the db for the test suite with the name of the new db
	s.DB, err = database.Init(connectionString, params, dbname, true)
	if err != nil {
		t.Fatalf("failed to init DB: %s", err)
	}
}

func (s *DatabaseSuite) Run(t *testing.T) {
	suite.Run(t, s)

	for _, dependentSuite := range s.suites {
		dependentSuite.Run(t)
	}
}

func (s *DatabaseSuite) AddDependentSuite(dependentSuite DependentSuite) {
	s.suites = append(s.suites, dependentSuite)
}
