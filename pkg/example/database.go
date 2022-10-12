package example

import (
	"database/sql"
	"embed"
	"net/http"

	"github.com/doug-martin/goqu/v9"
	"github.com/gorilla/mux"

	"github.com/charlieegan3/toolbelt/pkg/apis"
)

//go:embed database/migrations
var databaseToolMigrations embed.FS

// DatabaseTool is an example tool which demonstrates the use of the database feature
type DatabaseTool struct {
	db *sql.DB
}

func (d *DatabaseTool) Name() string {
	return "database"
}

func (d *DatabaseTool) FeatureSet() apis.FeatureSet {
	return apis.FeatureSet{
		HTTP:     true,
		Database: true,
		Jobs:     false,
	}
}

func (d *DatabaseTool) HTTPPath() string {
	return "database"
}

// SetConfig is a no-op for this tool
func (d *DatabaseTool) SetConfig(config map[string]any) error {
	return nil
}

func (d *DatabaseTool) DatabaseMigrations() (*embed.FS, string, error) {
	return &databaseToolMigrations, "database/migrations", nil
}

func (d *DatabaseTool) DatabaseSet(db *sql.DB) {
	d.db = db
}

func (d *DatabaseTool) Jobs() ([]apis.Job, error) {
	return []apis.Job{}, nil
}
func (d *DatabaseTool) ExternalJobsFuncSet(func(job apis.ExternalJob) error) {}

func (d *DatabaseTool) HTTPAttach(router *mux.Router) error {
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		goquDB := goqu.New("postgres", d.db)

		sel := goquDB.From("databasetool.example").
			Select("note").
			Where(goqu.Ex{
				"note": "database value",
			})

		var notes []struct {
			Note string `db:"note"`
		}

		if err := sel.ScanStructs(&notes); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(notes) != 1 {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err := writer.Write([]byte(notes[0].Note))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	})

	return nil
}
