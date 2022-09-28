package example

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/charlieegan3/toolbelt/pkg/apis"
	utilshttp "github.com/charlieegan3/toolbelt/pkg/utils/http"
	"github.com/gorilla/mux"
	"net/http"
)

// ConfigTool is an example tool which uses a configuration value in it's HTTP handler
type ConfigTool struct {
	exampleValue string
}

func (c *ConfigTool) Name() string {
	return "config-tool"
}

func (c *ConfigTool) FeatureSet() apis.FeatureSet {
	// this tool makes use of the config feature
	return apis.FeatureSet{
		HTTP:   true,
		Config: true,
	}
}

func (c *ConfigTool) HTTPPath() string {
	return "example-config-tool"
}

func (c *ConfigTool) SetConfig(config map[string]any) error {
	val, ok := config["exampleValue"].(string)
	if !ok {
		return fmt.Errorf("failed to get exampleValue from config")
	}

	c.exampleValue = val

	return nil
}

func (c *ConfigTool) DatabaseMigrations() (*embed.FS, string, error) {
	return nil, "", fmt.Errorf("not implemented")
}

func (c *ConfigTool) DatabaseSet(db *sql.DB) {}

func (c *ConfigTool) HTTPAttach(router *mux.Router) error {
	router.HandleFunc("", utilshttp.BuildRedirectHandler(c.HTTPPath()+"/")).Methods("GET")

	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte(c.exampleValue))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	})

	return nil
}
