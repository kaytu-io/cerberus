// server/web.go
package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kaytu-io/cerberus/workflow"
)

// StartServer initializes and starts the web server
func StartServer(wm *workflow.WorkflowManager) {
	http.HandleFunc("/status", HandleStatus)
	http.HandleFunc("/run/", func(w http.ResponseWriter, r *http.Request) {
		HandleRunWorkflow(w, r, wm)
	})

	// Start the HTTP server on port 8080
	http.ListenAndServe(":8080", nil)
}

// HandleStatus handles the /status endpoint
func HandleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service is running")
}

// HandleRunWorkflow triggers the specified workflow
func HandleRunWorkflow(w http.ResponseWriter, r *http.Request, wm *workflow.WorkflowManager) {
	// Extract workflow name from URL
	path := strings.TrimPrefix(r.URL.Path, "/run/")
	workflowName := strings.TrimSpace(path)
	if workflowName == "" {
		workflowName = workflow.DefaultWorkflow
	}

	err := wm.RunWorkflow(workflowName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error running workflow %s: %v", workflowName, err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Workflow %s executed successfully", workflowName)
}
