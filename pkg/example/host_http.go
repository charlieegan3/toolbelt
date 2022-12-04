package example

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/charlieegan3/toolbelt/pkg/apis"
	"github.com/gorilla/mux"
	"net/http"
)

// HostHTTPTool is a tool which doesn't use an http prefix, but rather a host matcher
type HostHTTPTool struct {
}

func (h *HostHTTPTool) Name() string {
	return "host-http-tool"
}

func (h *HostHTTPTool) FeatureSet() apis.FeatureSet {
	return apis.FeatureSet{
		HTTP:     true,
		HTTPHost: true,
	}
}
func (h *HostHTTPTool) HTTPPath() string {
	// unused
	return ""
}

func (h *HostHTTPTool) HTTPHost() string {
	return "example.com"
}

func (h *HostHTTPTool) HTTPAttach(router *mux.Router) error {
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("host tool"))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	})

	return nil
}

func (h *HostHTTPTool) DatabaseMigrations() (*embed.FS, string, error) {
	return nil, "", fmt.Errorf("not implemented")
}
func (h *HostHTTPTool) SetConfig(config map[string]any) error                { return nil }
func (h *HostHTTPTool) DatabaseSet(db *sql.DB)                               {}
func (h *HostHTTPTool) Jobs() ([]apis.Job, error)                            { return []apis.Job{}, nil }
func (h *HostHTTPTool) ExternalJobsFuncSet(func(job apis.ExternalJob) error) {}
